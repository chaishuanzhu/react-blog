import { message } from 'antd';
import copy from 'copy-to-clipboard';
import React from 'react';

import { siteConfig } from '@/site.config';

import CopyIcon from './CopyIcon';
import CopyrightIcon from './CopyrightIcon';
import s from './index.scss';

interface Props {
  id: number;
  title: string;
}

const CopyRight: React.FC<Props> = ({ id, title }) => {
  const url = `${window.location.origin}/post/${id}`;

  const copyUrl = async () => {
    if (await copy(url)) {
      message.success('复制成功!');
    }
  };

  return (
    <div className={s.copyrightBox}>
      <CopyrightIcon className={s.copyrightIcon} />
      <div className={s.title}>{title}</div>
      <div className={s.urlBox}>
        <div className={s.url}>{url}</div>
        <div className={s.copyBtn} onClick={copyUrl}>
          <CopyIcon className={s.copyIcon} />
        </div>
      </div>
      <div className={s.text}>
        本站所有文章除特别声明外，均采用
        <a
          href='https://creativecommons.org/licenses/by-nc-sa/4.0/'
          target='_blank'
          className={s.copyrightName}
          rel='noreferrer'
        >
          CC BY-NC-SA 4.0
        </a>
        许可协议，转载请注明来自
        <a href={window.location.origin} target='_blank' className={s.copyrightName} rel='noreferrer'>
          {siteConfig.title}
        </a>
        。
      </div>
    </div>
  );
};

export default CopyRight;
