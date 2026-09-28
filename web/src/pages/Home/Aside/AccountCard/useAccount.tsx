import { GithubOutlined, QqOutlined, WechatOutlined } from '@ant-design/icons';
import React from 'react';

import { siteConfig } from '@/site.config';

import Csdn from './Csdn';

export const useAccount = () => {
  const imgStyle = { width: '120px', height: '120px' };
  const { github, csdn, weChatQRCode, qqQRCode } = siteConfig.social;

  return [
    {
      show: !!github,
      isLink: true,
      link: github,
      ico: <GithubOutlined />,
      content: null
    },
    {
      show: !!csdn,
      isLink: true,
      link: csdn,
      ico: <Csdn />,
      content: null
    },
    {
      show: !!weChatQRCode,
      isLink: false,
      link: '',
      ico: <WechatOutlined />,
      content: <img src={weChatQRCode} style={imgStyle} />
    },
    {
      show: !!qqQRCode,
      isLink: false,
      link: '',
      ico: <QqOutlined />,
      content: <img src={qqQRCode} style={imgStyle} />
    }
  ].filter(item => item.show);
};
