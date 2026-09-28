import { useRequest, useSafeState } from 'ahooks';
import React from 'react';
import { useNavigate } from 'react-router';

import MyPagination from '@/components/MyPagination';
import { getArticles } from '@/utils/api';
import { homeSize, staleTime } from '@/utils/constant';

import s from './index.scss';
import PostCard from './PostCard';

const Section: React.FC = () => {
  const navigate = useNavigate();
  const [page, setPage] = useSafeState(1);

  const { data, loading } = useRequest(() => getArticles({ page, pageSize: homeSize }), {
    retryCount: 3,
    refreshDeps: [page],
    cacheKey: `Section-articles-${page}`,
    staleTime
  });

  return (
    <section className={s.section}>
      {data?.items.map(({ id, title, summary, publishedAt, tags }) => (
        <PostCard
          key={id}
          title={title}
          content={summary}
          date={publishedAt}
          tags={tags.map(tag => tag.name)}
          loading={loading}
          onClick={() => navigate(`/post/${id}`)}
        />
      ))}
      <MyPagination
        current={page}
        defaultPageSize={homeSize}
        total={data?.total}
        setPage={setPage}
        autoScroll={true}
        scrollToTop={document.body.clientHeight - 80}
      />
    </section>
  );
};

export default Section;
