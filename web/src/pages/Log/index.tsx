import { useRequest } from 'ahooks';
import React from 'react';

import Layout from '@/components/Layout';
import { getChangelogs } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import TimeItem from './TimeItem';

const Log: React.FC = () => {
  const { data, loading } = useRequest(getChangelogs, {
    retryCount: 3,
    cacheKey: 'changelogs',
    staleTime
  });

  return (
    <Layout title={Title.Log} loading={loading}>
      {data?.map(({ id, loggedAt, items }) => (
        <TimeItem key={id} date={loggedAt} logContent={items} />
      ))}
    </Layout>
  );
};

export default Log;
