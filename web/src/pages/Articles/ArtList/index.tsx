import dayjs from 'dayjs';
import React from 'react';
import { useNavigate } from 'react-router';

import DisplayBar from '@/components/DisplayBar';
import type { ArticleType } from '@/pages/constant';

import s from './index.scss';

interface Props {
  articles?: ArticleType[];
  loading?: boolean;
}

const ArtList: React.FC<Props> = ({ articles, loading }) => {
  const navigate = useNavigate();

  return (
    <>
      {articles?.length ? (
        articles?.map((item: ArticleType) => (
          <DisplayBar
            key={item.id}
            content={item.title}
            right={dayjs(item.publishedAt).format('YYYY-MM-DD')}
            onClick={() => navigate(`/post/${item.id}`)}
            loading={loading}
          />
        ))
      ) : (
        <div className={s.none}>暂时无相应文章 ~</div>
      )}
    </>
  );
};

export default ArtList;
