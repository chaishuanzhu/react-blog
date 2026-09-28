import type { CategoryList } from '@/utils/api';

const getChartData = (categories?: CategoryList) => {
  if (!categories) return [];
  const res = categories.items.map(item => ({ name: item.name, value: item.articleCount }));
  if (categories.uncategorizedCount) {
    res.push({ name: '未分类', value: categories.uncategorizedCount });
  }
  return res;
};

export const useOption = (categories: CategoryList | undefined, mode: number) => {
  const data = getChartData(categories);

  const labelColor = ['rgb(255, 255, 255)', 'rgb(53, 53, 53)', 'rgb(53, 53, 53)'];
  const backgroundColor = ['rgb(22, 54, 51)', 'rgb(157, 222, 255)', 'rgb(194, 209, 223)'];

  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: backgroundColor[mode],
      borderColor: backgroundColor[mode],
      textStyle: {
        color: labelColor[mode],
        fontSize: 16,
        fontFamily: 'dengxian'
      }
    },
    series: [
      {
        type: 'pie',
        radius: '88%',
        height: '400px',
        data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          color: labelColor[mode],
          fontSize: 18,
          fontFamily: 'dengxian'
        }
      }
    ]
  };
};
