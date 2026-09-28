import './index.custom.scss';

import { Button, Input, Select } from '@arco-design/web-react';
import { useRequest, useTitle } from 'ahooks';
import React, { useState } from 'react';
import { BiBrushAlt, BiSearch } from 'react-icons/bi';
import { useNavigate, useSearchParams } from 'react-router';

import MyTable from '@/components/MyTable';
import PageHeader from '@/components/PageHeader';
import { categoryApi, tagApi } from '@/utils/api';
import { siteTitle } from '@/utils/constant';

import { Title } from '../titleConfig';
import s from './index.scss';
import { useArticleTable } from './useArticleTable';

const toId = (raw: string | null) => (raw ? Number(raw) || undefined : undefined);

const Article: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Articles}`);
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const keyword = searchParams.get('keyword') || undefined;
  const categoryId = toId(searchParams.get('categoryId'));
  const tagId = toId(searchParams.get('tagId'));
  const [keywordInput, setKeywordInput] = useState(keyword ?? '');

  const { data: categories, loading: classLoading } = useRequest(categoryApi.list);
  const { data: tags = [], loading: tagLoading } = useRequest(tagApi.list);

  const { columns, data, total, loading, page, setPage } = useArticleTable('published', {
    keyword,
    categoryId,
    tagId
  });

  const setFilter = (key: string, value?: string | number) => {
    setSearchParams(prev => {
      const params = new URLSearchParams(prev);
      if (value === undefined || value === '') params.delete(key);
      else params.set(key, String(value));
      params.set('page', '1');
      return params;
    });
  };

  const clearSearch = () => {
    setKeywordInput('');
    setSearchParams({ page: '1' });
  };

  return (
    <>
      <PageHeader text='写文章' onClick={() => navigate('/addArticle')}>
        <div className={s.searchBox}>
          <div className={s.search}>
            <Input
              size='large'
              allowClear
              style={{ flex: 1, marginRight: 10 }}
              className='articleInputBox'
              placeholder='输入文章标题关键字'
              value={keywordInput}
              onChange={value => setKeywordInput(value)}
              onPressEnter={() => setFilter('keyword', keywordInput.trim())}
              onClear={() => setFilter('keyword')}
            />
            <Select
              size='large'
              placeholder='请选择文章分类'
              style={{ flex: 1, marginRight: 10 }}
              showSearch
              allowClear
              value={categoryId}
              onChange={value => setFilter('categoryId', value)}
              disabled={classLoading}
              options={(categories?.items ?? []).map(({ id, name }) => ({
                value: id,
                label: name
              }))}
            />
            <Select
              placeholder='请选择文章标签'
              size='large'
              style={{ flex: 1, marginRight: 10 }}
              showSearch
              allowClear
              value={tagId}
              onChange={value => setFilter('tagId', value)}
              disabled={tagLoading}
              options={tags.map(({ id, name }) => ({ value: id, label: name }))}
            />
          </div>
          <div>
            <Button
              type='primary'
              size='large'
              onClick={() => setFilter('keyword', keywordInput.trim())}
              style={{ fontSize: 16, marginRight: 10 }}
            >
              <BiSearch />
            </Button>
            <Button type='primary' size='large' onClick={clearSearch} style={{ fontSize: 16 }}>
              <BiBrushAlt />
            </Button>
          </div>
        </div>
      </PageHeader>
      <MyTable
        loading={loading}
        columns={columns}
        data={data}
        total={total}
        page={page}
        setPage={setPage}
      />
    </>
  );
};

export default Article;
