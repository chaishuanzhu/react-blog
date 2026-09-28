import { Button, Message } from '@arco-design/web-react';
import { useRequest, useTitle } from 'ahooks';
import classNames from 'classnames';
import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router';

import MarkDown from '@/components/MarkDown';
import PageHeader from '@/components/PageHeader';
import { pageApi } from '@/utils/api';
import { siteTitle } from '@/utils/constant';
import { mutate } from '@/utils/feedback';
import { useScrollSync } from '@/utils/hooks/useScrollSync';

import { Title } from '../titleConfig';
import s from './index.scss';

const AboutEdit: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const isMe = searchParams.get('me') === '1';
  const key = isMe ? 'about-me' : 'about-site';

  useTitle(`${siteTitle} | ${isMe ? Title.AboutMe : Title.AboutSite}`);

  const { leftRef, rightRef, handleScrollRun } = useScrollSync();
  const [content, setContent] = useState('');

  useRequest(() => pageApi.get(key), { onSuccess: setContent });

  const updateAbout = async () => {
    if (!content.trim()) {
      Message.info('请写点什么再更新！');
      return;
    }
    if (await mutate(() => pageApi.update(key, content), '更新成功！')) {
      navigate('/about');
    }
  };

  return (
    <>
      <PageHeader text='返回' onClick={() => navigate('/about')}>
        <div className={s.aboutTitle}>关于{isMe ? '我' : '本站'}</div>
        <Button size='large' type='primary' className={s.aboutUpdate} onClick={updateAbout}>
          更新
        </Button>
      </PageHeader>
      <div className={s.markedEditBox}>
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

export default AboutEdit;
