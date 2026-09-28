import { useRequest } from 'ahooks';
import React from 'react';

import Card from '@/components/Card';
import { getTags } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import s from './index.scss';

const TagCard: React.FC = () => {
  const { data, loading } = useRequest(getTags, {
    retryCount: 3,
    cacheKey: 'tags',
    staleTime
  });

  return (
    <Card className={s.card} loading={loading}>
      {data?.map(item => (
        <span className={s.tag} key={item.id}>
          {item.name}
        </span>
      ))}
    </Card>
  );
};

export default TagCard;
