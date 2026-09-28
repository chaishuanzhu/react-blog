import { PieChart } from 'echarts/charts';
import { LegendComponent, TitleComponent, TooltipComponent } from 'echarts/components';
import * as echarts from 'echarts/core';
import { LabelLayout } from 'echarts/features';
import { CanvasRenderer } from 'echarts/renderers';
import ReactEChartsCore from 'echarts-for-react/lib/core';
import React from 'react';
import { connect } from 'react-redux';

import type { storeState } from '@/redux/interface';
import type { CategoryList } from '@/utils/api';

import s from './index.scss';
import { useOption } from './useOption';

interface Props {
  categories?: CategoryList;
  mode?: number;
}

echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  PieChart,
  CanvasRenderer,
  LabelLayout
]);

const Chart: React.FC<Props> = ({ categories, mode }) => {
  const option = useOption(categories, mode!);

  return (
    <div className={s.box}>
      <h3>📊文章分布</h3>
      <ReactEChartsCore
        style={{
          height: '400px'
        }}
        echarts={echarts}
        option={option}
        notMerge={true}
        lazyUpdate={true}
        theme='theme_name'
      />
    </div>
  );
};

export default connect((state: storeState) => ({
  mode: state.mode
}))(Chart);
