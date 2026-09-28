import { IconLoading } from '@arco-design/web-react/icon';
import classNames from 'classnames';
import { PieChart } from 'echarts/charts';
import { LegendComponent, TitleComponent, TooltipComponent } from 'echarts/components';
import * as echarts from 'echarts/core';
import { LabelLayout } from 'echarts/features';
import { CanvasRenderer } from 'echarts/renderers';
import ReactEChartsCore from 'echarts-for-react/lib/core';
import React from 'react';
import { useNavigate } from 'react-router';

import type { CategoryList } from '@/utils/api';

import s from './index.scss';

echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  PieChart,
  CanvasRenderer,
  LabelLayout
]);

interface Props {
  categories?: CategoryList;
  loading: boolean;
}

const ChartCard: React.FC<Props> = ({ categories, loading }) => {
  const navigate = useNavigate();

  const data = [
    ...(categories?.items ?? [])
      .filter(c => c.articleCount > 0)
      .map(c => ({ id: c.id, name: c.name, value: c.articleCount })),
    ...(categories?.uncategorizedCount
      ? [{ id: 0, name: '未分类', value: categories.uncategorizedCount }]
      : [])
  ];

  const option = {
    tooltip: {
      trigger: 'item',
      textStyle: { fontSize: 16, fontFamily: 'dengxian' }
    },
    series: [
      {
        type: 'pie',
        radius: '80%',
        height: '100%',
        data,
        emphasis: {
          itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' }
        },
        label: { fontSize: 18, fontFamily: 'dengxian' }
      }
    ]
  };

  const onEvents = {
    click: (params: { data: { id: number } }) => {
      if (params.data.id) navigate(`/article?categoryId=${params.data.id}`);
    }
  };

  return (
    <div className={classNames(s.chartBox, { [s.loadingCenter]: loading })}>
      <div className={s.chartTitle}>文章概览</div>
      {loading ? (
        <IconLoading className={s.loading} />
      ) : (
        <ReactEChartsCore
          style={{ height: '100%' }}
          echarts={echarts}
          option={option}
          notMerge={true}
          lazyUpdate={true}
          onEvents={onEvents}
        />
      )}
    </div>
  );
};

export default ChartCard;
