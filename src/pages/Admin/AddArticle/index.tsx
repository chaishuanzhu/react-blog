import { Button, Input, Message, Select } from '@arco-design/web-react';
import { useRequest, useTitle } from 'ahooks';
import classNames from 'classnames';
import dayjs from 'dayjs';
import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router';

import MarkDown from '@/components/MarkDown';
import UploadButton from '@/components/UploadButton';
import type { ArticleStatus } from '@/utils/api';
import { articleApi, categoryApi, tagApi } from '@/utils/api';
import { dateTimeFormat, siteTitle } from '@/utils/constant';
import { mutate, parseLocalTime } from '@/utils/feedback';
import { useScrollSync } from '@/utils/hooks/useScrollSync';

import { Title } from '../titleConfig';
import s from './index.scss';

const AddArticle: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const id = Number(searchParams.get('id')) || 0;

  useTitle(`${siteTitle} | ${id ? Title.EditArticle : Title.AddArticle}`);

  const { leftRef, rightRef, handleScrollRun } = useScrollSync();

  const [title, setTitle] = useState('');
  const [categoryId, setCategoryId] = useState<number | undefined>();
  const [tagIds, setTagIds] = useState<number[]>([]);
  const [date, setDate] = useState(dayjs().format(dateTimeFormat));
  const [content, setContent] = useState('');
  const [status, setStatus] = useState<ArticleStatus>('draft');
  const [saving, setSaving] = useState(false);

  useRequest(() => articleApi.get(id), {
    ready: !!id,
    onSuccess: article => {
      setTitle(article.title);
      setCategoryId(article.category?.id);
      setTagIds(article.tags.map(t => t.id));
      setDate(dayjs(article.publishedAt).format(dateTimeFormat));
      setContent(article.content ?? '');
      setStatus(article.status);
    }
  });

  const { data: categories, loading: classLoading } = useRequest(categoryApi.list);
  const { data: tags = [], loading: tagLoading } = useRequest(tagApi.list);

  const insertAtCursor = (text: string) => {
    const el = leftRef.current;
    const start = el?.selectionStart ?? content.length;
    const end = el?.selectionEnd ?? content.length;
    setContent(content.slice(0, start) + text + content.slice(end));
  };

  const postArticle = async (next: ArticleStatus) => {
    if (!title.trim() || !date.trim()) {
      Message.info('请至少输入标题、时间！');
      return;
    }
    if (next === 'published' && !content.trim()) {
      Message.info('发布前请填写正文！');
      return;
    }
    const publishedAt = parseLocalTime(date, dateTimeFormat);
    if (!publishedAt) {
      Message.info('日期字符串不合法！');
      return;
    }

    const input = {
      title: title.trim(),
      content,
      categoryId: categoryId ?? null,
      tagIds,
      status: next,
      publishedAt
    };
    const successText = next === 'published' ? `${id ? '更新' : '发布'}文章成功！` : '保存草稿成功！';

    setSaving(true);
    const ok = await mutate(
      () => (id ? articleApi.update(id, input) : articleApi.create(input)),
      successText
    );
    setSaving(false);
    if (ok) navigate(next === 'published' ? '/article' : '/draft');
  };

  return (
    <>
      <div className={s.addArticleHeader}>
        <div className={s.top}>
          <Input
            className={s.chineseTitle}
            addBefore='标题'
            allowClear
            size='large'
            value={title}
            onChange={value => setTitle(value)}
          />
          <Button
            size='large'
            type='primary'
            style={{ marginRight: 10 }}
            loading={saving}
            onClick={() => postArticle('draft')}
          >
            存为草稿
          </Button>
          <Button
            size='large'
            type='primary'
            status='success'
            loading={saving}
            onClick={() => postArticle('published')}
          >
            {id && status === 'published' ? '更新' : '发布'}文章
          </Button>
        </div>
        <div className={s.bottom}>
          <Select
            addBefore='分类'
            size='large'
            className={s.classText}
            showSearch
            allowClear
            value={categoryId}
            onChange={value => setCategoryId(value)}
            disabled={classLoading}
            options={(categories?.items ?? []).map(({ id, name }) => ({
              value: id,
              label: name
            }))}
          />
          <Select
            addBefore='标签'
            size='large'
            className={s.tags}
            maxTagCount={6}
            mode='multiple'
            showSearch
            allowClear
            value={tagIds}
            onChange={value => setTagIds(value)}
            disabled={tagLoading}
            options={tags.map(({ id, name }) => ({ value: id, label: name }))}
          />
          <Input
            addBefore='时间'
            value={date}
            placeholder={dateTimeFormat}
            onChange={value => setDate(value)}
            className={s.time}
            allowClear
            size='large'
            style={{ marginRight: 10 }}
          />
          <UploadButton
            size='large'
            onUploaded={(url, file) => insertAtCursor(`![${file.name}](${url})`)}
          />
        </div>
      </div>
      <div className={s.contentEdit}>
        <textarea
          ref={leftRef}
          className={classNames(s.markedEdit, s.input)}
          value={content}
          onChange={e => setContent(e.target.value)}
          onScroll={handleScrollRun}
        />
        <MarkDown
          ref={rightRef}
          className={s.markedEdit}
          content={content}
          onScroll={handleScrollRun}
        />
      </div>
    </>
  );
};

export default AddArticle;
