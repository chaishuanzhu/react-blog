import { useMount, useSafeState, useTitle } from 'ahooks';
import { load } from 'jinrishici';
import React from 'react';
import { connect } from 'react-redux';

import PageTitle from '@/components/PageTitle';
import { setNavShow } from '@/redux/actions';
import { siteConfig } from '@/site.config';
import useTop from '@/utils/hooks/useTop';

import Aside from './Aside';
import s from './index.scss';
import Section from './Section';

interface Props {
  setNavShow?: Function;
}

const Home: React.FC<Props> = ({ setNavShow }) => {
  useTitle(siteConfig.title);
  useTop(setNavShow);

  const [poem, setPoem] = useSafeState('');
  useMount(() => {
    load(res => setPoem(res.data.content));
  });

  return (
    <>
      <PageTitle title={siteConfig.title} desc={poem || ''} className={s.homeTitle} />
      <div className={s.body}>
        <Section />
        <Aside />
      </div>
    </>
  );
};

export default connect(() => ({}), { setNavShow })(Home);
