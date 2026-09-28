import { useRequest, useSafeState } from 'ahooks';
import React from 'react';

import Layout from '@/components/Layout';
import MyPagination from '@/components/MyPagination';
import { getArticles } from '@/utils/api';
import { detailPostSize, staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import ArtList from './ArtList';
import Search from './Search';

const Articles: React.FC = () => {
  const [page, setPage] = useSafeState(1);
  const [keyword, setKeyword] = useSafeState('');

  const { data, loading } = useRequest(
    () => getArticles({ page, pageSize: detailPostSize, keyword: keyword || undefined }),
    {
      retryCount: 3,
      refreshDeps: [page, keyword],
      cacheKey: `Articles-${keyword}-${page}`,
      staleTime
    }
  );

  return (
    <Layout title={Title.Articles}>
      <Search page={page} setPage={setPage} keyword={keyword} setKeyword={setKeyword} />
      <ArtList articles={data?.items} loading={loading} />
      <MyPagination
        current={page}
        defaultPageSize={detailPostSize}
        total={data?.total}
        setPage={setPage}
        autoScroll={true}
        scrollToTop={440}
      />
    </Layout>
  );
};

export default Articles;
