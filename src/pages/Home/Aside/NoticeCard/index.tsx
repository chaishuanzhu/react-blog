import { useRequest } from 'ahooks';
import React from 'react';

import Card from '@/components/Card';
import { getSite } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import s from './index.scss';

const NoticeCard: React.FC = () => {
  const { data, loading } = useRequest(getSite, {
    retryCount: 3,
    cacheKey: 'site',
    staleTime
  });

  return (
    <Card loading={loading}>
      <div className={s.notice}>{data?.notice}</div>
    </Card>
  );
};

export default NoticeCard;
