import React from 'react';

import type { CategoryList } from '@/utils/api';

import AboutText from './AboutText';
import Chart from './Chart';

interface Props {
  content?: string;
  categories?: CategoryList;
  className?: string;
}

const AboutSite: React.FC<Props> = ({ content, categories, className }) => {
  return (
    <div className={className}>
      <Chart categories={categories} />
      <AboutText content={content} />
    </div>
  );
};

export default AboutSite;
