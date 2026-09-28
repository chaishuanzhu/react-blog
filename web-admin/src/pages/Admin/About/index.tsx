import { useRequest, useTitle } from 'ahooks';
import React from 'react';

import AboutBase from '@/components/AboutBase';
import { pageApi } from '@/utils/api';
import { siteTitle } from '@/utils/constant';

import { Title } from '../titleConfig';
import s from './index.scss';

const About: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.About}`);

  const { data: aboutMe } = useRequest(() => pageApi.get('about-me'));
  const { data: aboutSite } = useRequest(() => pageApi.get('about-site'));

  return (
    <div className={s.aboutBox}>
      <div className={s.left}>
        <AboutBase content={aboutMe} site='关于我' params={1} />
      </div>
      <div className={s.right}>
        <AboutBase content={aboutSite} site='关于本站' params={0} />
      </div>
    </div>
  );
};

export default About;
