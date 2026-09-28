import { useRequest } from 'ahooks';
import React from 'react';
import { useNavigate } from 'react-router';

import Card from '@/components/Card';
import { getSite } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import s from './index.scss';

const DataCard: React.FC = () => {
  const navigate = useNavigate();
  const { data, loading } = useRequest(getSite, {
    retryCount: 3,
    cacheKey: 'site',
    staleTime
  });

  return (
    <Card className={s.card} loading={loading}>
      <div className={s.blogData} onClick={() => navigate('/articles')}>
        <div className={s.name}>文章</div>
        <div className={s.num}>{data?.articleCount}</div>
      </div>
      <div className={s.blogData} onClick={() => navigate('/classes')}>
        <div className={s.name}>分类</div>
        <div className={s.num}>{data?.categoryCount}</div>
      </div>
      <div className={s.blogData} onClick={() => navigate('/tags')}>
        <div className={s.name}>标签</div>
        <div className={s.num}>{data?.tagCount}</div>
      </div>
    </Card>
  );
};

export default DataCard;
