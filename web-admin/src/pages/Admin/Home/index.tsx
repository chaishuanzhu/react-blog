import { useRequest, useTitle } from 'ahooks';
import React from 'react';

import ChartCard from '@/components/ChartCard';
import ClassCard from '@/components/ClassCard';
import CountCard from '@/components/CountCard';
import NoticeCard from '@/components/NoticeCard';
import TagCard from '@/components/TagCard';
import type { Stats } from '@/utils/api';
import { categoryApi, getStats } from '@/utils/api';
import { siteTitle } from '@/utils/constant';

import s from './index.scss';

const countCards: { label: string; key: keyof Stats }[] = [
  { label: '文章数', key: 'publishedCount' },
  { label: '说说数', key: 'momentCount' },
  { label: '留言数', key: 'commentCount' },
  { label: '友链数', key: 'friendLinkCount' },
  { label: '访问量', key: 'viewCount' }
];

const Home: React.FC = () => {
  useTitle(siteTitle);

  const { data: stats, loading: statsLoading } = useRequest(getStats);
  const {
    data: categories,
    loading: categoriesLoading,
    refresh: refreshCategories
  } = useRequest(categoryApi.list);

  return (
    <>
      {/* 统计卡片区 */}
      <div className={s.countCardContainer}>
        {countCards.map(({ label, key }) => (
          <CountCard key={key} label={label} value={stats?.[key]} loading={statsLoading} />
        ))}
      </div>
      {/* 扇形图、分类、标签、公告 */}
      <div className={s.homeBigContainer}>
        <div className={s.chartContainer}>
          <ChartCard categories={categories} loading={categoriesLoading} />
        </div>
        <div className={s.classesContainer}>
          <ClassCard
            categories={categories}
            loading={categoriesLoading}
            onChanged={refreshCategories}
          />
        </div>
        <div className={s.tagsNoticeContainer}>
          <div className={s.NoticeContainer}>
            <NoticeCard />
          </div>
          <div className={s.tagsContainer}>
            <TagCard />
          </div>
        </div>
      </div>
    </>
  );
};

export default Home;
