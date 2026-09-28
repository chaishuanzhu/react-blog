import { IconLoading } from '@arco-design/web-react/icon';
import classNames from 'classnames';
import React from 'react';

import s from './index.scss';

interface Props {
  label: string;
  value?: number;
  loading: boolean;
  className?: string;
}

const CountCard: React.FC<Props> = ({ label, value, loading, className }) => {
  return (
    <div className={classNames(s.countCardBox, className)}>
      <div className={s.key}>{label}</div>
      <div className={classNames(s.value, { [s.loading]: loading })}>
        {loading ? <IconLoading /> : (value ?? 0)}
      </div>
    </div>
  );
};

export default CountCard;
