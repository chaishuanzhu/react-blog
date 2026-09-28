import { useRequest } from 'ahooks';
import React, { useState } from 'react';

import ImgView from '@/components/ImgView';
import Layout from '@/components/Layout';
import { getMoments } from '@/utils/api';
import { staleTime } from '@/utils/constant';

import { Title } from '../titleConfig';
import SayPop from './SayPop';

const Say: React.FC = () => {
  const { data, loading } = useRequest(() => getMoments(), {
    retryCount: 3,
    cacheKey: 'moments',
    staleTime
  });

  const [url, setUrl] = useState('');
  const [showPreView, setShowPreView] = useState(false);

  const handlePreView = (url: string) => {
    setShowPreView(true);
    setUrl(url);
  };

  return (
    <Layout title={Title.Say} loading={loading}>
      {data?.items.map(({ id, content, createdAt, images }) => (
        <SayPop
          key={id}
          content={content}
          date={createdAt}
          imgs={images}
          handlePreView={handlePreView}
        />
      ))}

      <ImgView
        viewUrl={url}
        isViewShow={showPreView}
        onClick={() => setShowPreView(false)}
      />
    </Layout>
  );
};

export default Say;
