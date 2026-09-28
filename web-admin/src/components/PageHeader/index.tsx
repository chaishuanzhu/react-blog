import { Button } from '@arco-design/web-react';
import React from 'react';

import s from './index.scss';

interface Props {
  text: string;
  onClick: () => void;
  children?: React.ReactNode;
}

const PageHeader: React.FC<Props> = ({ text, onClick, children }) => {
  return (
    <div className={s.pageHeaderBox}>
      <Button type='primary' size='large' onClick={onClick}>
        {text}
      </Button>
      {children}
    </div>
  );
};

export default PageHeader;
