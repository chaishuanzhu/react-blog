import { useRequest } from 'ahooks';
import React from 'react';

import Layout from '@/components/Layout';
import { getFriendLinks } from '@/utils/api';
import { staleTime } from '@/utils/constant';
import { shuffleArray } from '@/utils/function';

import { Title } from '../titleConfig';
import s from './index.scss';
import LinkItem from './LinkItem';

const Link: React.FC = () => {
  const { data, loading } = useRequest(getFriendLinks, {
    retryCount: 3,
    cacheKey: 'friend-links',
    staleTime
  });

  return (
    <Layout title={Title.Link} loading={loading} className={s.box}>
      {shuffleArray(data || []).map(item => (
        <LinkItem
          key={item.id}
          link={item.url}
          avatar={item.avatar}
          name={item.name}
          descr={item.description}
        />
      ))}
    </Layout>
  );
};

export default Link;
