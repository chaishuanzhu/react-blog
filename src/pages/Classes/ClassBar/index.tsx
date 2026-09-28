import classNames from 'classnames';
import type { MouseEventHandler } from 'react';
import React from 'react';

import s from './index.scss';

interface Props {
  content: string;
  num: number;
  className?: string;
  onClick?: MouseEventHandler<HTMLDivElement>;
}

const ClassBar: React.FC<Props> = ({ content, num, onClick, className }) => {
  return (
    <div className={classNames(s.classBar, className)} onClick={onClick}>
      {content}
      <div className={s.classNum}>{num}</div>
    </div>
  );
};

export default ClassBar;
