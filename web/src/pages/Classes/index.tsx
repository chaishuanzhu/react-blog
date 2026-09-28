import { useRequest } from 'ahooks';
import React from 'react';
import { useNavigate } from 'react-router';

import Layout from '@/components/Layout';
import { getCategories } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import ClassBar from './ClassBar';
import s from './index.scss';

const Classes: React.FC = () => {
  const navigate = useNavigate();

  const { data, loading } = useRequest(getCategories, {
    retryCount: 3,
    cacheKey: 'categories',
    staleTime
  });

  return (
    <Layout title={Title.Classes} loading={loading} className={s.classBox} rows={8}>
      {data?.items.map(item => (
        <ClassBar
          className={s.classItem}
          key={item.id}
          content={item.name}
          num={item.articleCount}
          onClick={() => navigate(`/artDetail?class=${encodeURIComponent(item.name)}`)}
        />
      ))}
    </Layout>
  );
};

export default Classes;
