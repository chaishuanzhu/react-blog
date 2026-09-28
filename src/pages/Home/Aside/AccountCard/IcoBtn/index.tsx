import { Popover } from 'antd';
import type { ReactNode } from 'react';
import React from 'react';

import s from './index.scss';

interface Props {
  isLink: boolean;
  link?: string;
  content?: ReactNode;
  children?: ReactNode;
}

const IcoBtn: React.FC<Props> = ({ isLink, link, content, children }) => {
  return isLink ? (
    <a className={s.socialBtn} href={link} target='_blank' rel='noreferrer'>
      {children}
    </a>
  ) : (
    <Popover trigger='hover' content={content} classNames={{ root: s.card }}>
      <div className={s.socialBtn}>{children}</div>
    </Popover>
  );
};

export default IcoBtn;
