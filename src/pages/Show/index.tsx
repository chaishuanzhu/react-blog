import { useRequest } from 'ahooks';
import React from 'react';

import Layout from '@/components/Layout';
import { getProjects } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import s from './index.scss';
import ShowItem from './ShowItem';

const Show: React.FC = () => {
  const { data, loading } = useRequest(getProjects, {
    retryCount: 3,
    cacheKey: 'projects',
    staleTime
  });

  return (
    <Layout title={Title.Show} loading={loading} className={s.showBox}>
      {data?.map(item => (
        <ShowItem
          key={item.id}
          cover={item.cover}
          link={item.url}
          name={item.name}
          descr={item.description}
        />
      ))}
    </Layout>
  );
};

export default Show;
