import { Skeleton } from 'antd';
import classNames from 'classnames';
import type { MouseEventHandler } from 'react';
import React from 'react';

import s from './index.scss';

interface Props {
  className?: string;
  loading?: boolean;
  isStatic?: boolean;
  onClick?: MouseEventHandler<HTMLDivElement>;
  children?: React.ReactNode;
}

const Card: React.FC<Props> = ({ children, className, loading, isStatic, onClick }) => {
  return (
    <div
      className={classNames(
        s.card,
        { [s.center]: loading },
        { [s.active]: !isStatic },
        className
      )}
      onClick={onClick}
    >
      {loading ? <Skeleton paragraph={{ rows: 1 }} /> : children}
    </div>
  );
};

export default Card;
