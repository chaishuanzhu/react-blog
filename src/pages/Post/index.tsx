import { useRequest } from 'ahooks';
import React from 'react';
import { useParams } from 'react-router';

import Comment from '@/components/Comment';
import Layout from '@/components/Layout';
import MarkDown from '@/components/MarkDown';
import { getArticle } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import CopyRight from './CopyRight';
import s from './index.scss';
import Navbar from './Navbar';
import PostTags from './PostTags';

const Post: React.FC = () => {
  const { id = '' } = useParams();

  const { data, loading, error } = useRequest(() => getArticle(id), {
    refreshDeps: [id],
    cacheKey: `Post-${id}`,
    staleTime
  });

  return (
    <Layout
      title={error && !data ? '文章不存在' : data?.title}
      loading={loading && !data}
      classes={data?.category?.name}
      date={data?.publishedAt}
      isPost={true}
      rows={14}
    >
      {data && (
        <>
          <div id='post-content'>
            <MarkDown content={data.content} className={s.mb} />
          </div>
          <PostTags tags={data.tags.map(tag => tag.name)} />
          <CopyRight id={data.id} title={data.title} />
          <Comment articleId={data.id} />
          <Navbar content={data.content} />
        </>
      )}
    </Layout>
  );
};

export default Post;
