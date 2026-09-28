import './index.custom.scss';

import { VerticalAlignTopOutlined } from '@ant-design/icons';
import { useScroll } from 'ahooks';
import React from 'react';
import { connect } from 'react-redux';

import { setNavShow } from '@/redux/actions';

import s from './index.scss';

interface Props {
  setNavShow?: Function;
}

const visibilityHeight = 300;

const BackToTop: React.FC<Props> = ({ setNavShow }) => {
  const scroll = useScroll(document);

  const backTop = () => {
    setNavShow?.(true);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if ((scroll?.top ?? 0) < visibilityHeight) return null;

  return (
    <div className={`BackTop ${s.box}`} onClick={backTop}>
      <div className={s.backTop}>
        <VerticalAlignTopOutlined />
      </div>
    </div>
  );
};

export default connect(() => ({}), {
  setNavShow
})(BackToTop);
