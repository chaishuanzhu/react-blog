import { useRequest } from 'ahooks';
import React from 'react';
import { useNavigate } from 'react-router';

import Layout from '@/components/Layout';
import { getTags } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import s from './index.scss';

const Tags: React.FC = () => {
  const navigate = useNavigate();

  const { data, loading } = useRequest(getTags, {
    retryCount: 3,
    cacheKey: 'tags',
    staleTime
  });

  return (
    <Layout title={Title.Tags} loading={loading} className={s.tagsBox} rows={3}>
      {data?.map(item => (
        <span
          className={s.tagItem}
          key={item.id}
          onClick={() => navigate(`/artDetail?tag=${encodeURIComponent(item.name)}`)}
        >
          {item.name}
        </span>
      ))}
    </Layout>
  );
};

export default Tags;
