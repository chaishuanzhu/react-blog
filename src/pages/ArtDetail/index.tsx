import { useRequest, useSafeState } from 'ahooks';
import dayjs from 'dayjs';
import React from 'react';
import { useNavigate, useSearchParams } from 'react-router';

import DisplayBar from '@/components/DisplayBar';
import Layout from '@/components/Layout';
import MyPagination from '@/components/MyPagination';
import { getArticles } from '@/utils/api';
import { detailPostSize, staleTime } from '@/utils/constant';

const ArtDetail: React.FC = () => {
  const [searchParams] = useSearchParams();
  const tag = searchParams.get('tag') || '';
  const category = searchParams.get('class') || '';
  const navigate = useNavigate();

  const [page, setPage] = useSafeState(1);

  const filter = tag ? { tag } : { category };

  const { data, loading } = useRequest(
    () => getArticles({ ...filter, page, pageSize: detailPostSize }),
    {
      retryCount: 3,
      refreshDeps: [page, tag, category],
      cacheKey: `ArtDetail-${JSON.stringify(filter)}-${page}`,
      staleTime
    }
  );

  return (
    <Layout title={tag || category}>
      {data?.items.map(item => (
        <DisplayBar
          key={item.id}
          content={item.title}
          right={dayjs(item.publishedAt).format('YYYY-MM-DD')}
          loading={loading}
          onClick={() => navigate(`/post/${item.id}`)}
        />
      ))}
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

export default ArtDetail;
