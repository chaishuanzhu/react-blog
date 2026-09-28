import React from 'react';
import { useNavigate } from 'react-router';

import MarkDown from '../MarkDown';
import PageHeader from '../PageHeader';
import s from './index.scss';

interface Props {
  content?: string;
  site: string;
  params: 0 | 1;
}

const AboutBase: React.FC<Props> = ({ content = '', site, params }) => {
  const navigate = useNavigate();

  return (
    <>
      <PageHeader text='编辑' onClick={() => navigate(`/aboutEdit?me=${params}`)}>
        <div className={s.site}>{site}</div>
      </PageHeader>
      <div className={s.markDownContent}>
        <MarkDown content={content} className={s.contentBox} />
      </div>
    </>
  );
};

export default AboutBase;
