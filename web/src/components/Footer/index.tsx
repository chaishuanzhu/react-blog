import React from 'react';

import { siteConfig } from '@/site.config';

import s from './index.scss';

const Footer: React.FC = () => {
  const frameArr = [
    'React',
    'React Router',
    'Redux',
    'Webpack',
    'AntD',
    'ahooks',
    'Go',
    'MySQL'
  ];
  const { sourceUrl, icp, police } = siteConfig;

  return (
    <footer className={s.footer}>
      <span>
        个人博客系统
        {sourceUrl && (
          <a href={sourceUrl} target='_blank' rel='noreferrer' className={s.text}>
            「源代码」
          </a>
        )}
      </span>
      {icp.no && (
        <span>
          <a href={icp.url} target='_blank' rel='noreferrer' className={s.text}>
            {icp.no}
          </a>
        </span>
      )}
      {police.no && (
        <span>
          <a href={police.url} target='_blank' rel='noreferrer' className={s.text}>
            {police.no}
          </a>
        </span>
      )}
      <span>
        {frameArr.map((item, index) => (
          <span className={s.siteFrame} key={index}>
            {item}
          </span>
        ))}
      </span>
    </footer>
  );
};

export default Footer;
