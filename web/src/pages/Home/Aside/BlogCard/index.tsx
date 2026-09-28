import React from 'react';

import Card from '@/components/Card';
import { siteConfig } from '@/site.config';
import { useTime } from '@/utils/hooks/useTime';

import s from './index.scss';

const BlogCard: React.FC = () => {
  const { timeText } = useTime();

  return (
    <Card className={s.card}>
      <p className={s.text}>
        {timeText}，<br />
        我叫<span className={s.color}>{siteConfig.author.name}</span>，<br />
        欢迎来到
        <br />
        我的<span className={s.color}>个人博客</span>。
      </p>
      <img src={siteConfig.cardImage} className={s.avatar} />
    </Card>
  );
};

export default BlogCard;
