import './index.custom.scss';

import { MenuFoldOutlined } from '@ant-design/icons';
import { useBoolean } from 'ahooks';
import { Drawer } from 'antd';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';

import { setNavShow } from '@/redux/actions';

import s from './index.scss';
import Toc from './Toc';

export const POST_CONTENT = '#post-content';

interface Props {
  content?: string;
  setNavShow?: Function;
}

const Navbar: React.FC<Props> = ({ content, setNavShow }) => {
  const [visible, { setTrue: openDrawer, setFalse: closeDrawer }] = useBoolean(false);

  return (
    <>
      {/* 正常的目录 */}
      <Toc
        target={POST_CONTENT}
        content={content}
        className={classNames('postNavBar', s.navBar)}
        headingTopOffset={15}
        onNavItemClick={() => setNavShow?.(false)}
      />
      {/* 中屏显示的按钮 */}
      <div className={s.hoverBar} onClick={openDrawer}>
        <MenuFoldOutlined />
      </div>
      {/* 中屏抽屉 */}
      <Drawer
        placement='right'
        onClose={closeDrawer}
        open={visible}
        rootClassName={classNames(s.drawer, 'mobile-navBar-box')}
        size={340}
      >
        <Toc
          target={POST_CONTENT}
          content={content}
          className='postNavBar'
          headingTopOffset={15 + 60}
          onNavItemClick={() => setNavShow?.(true)}
        />
      </Drawer>
    </>
  );
};

export default connect(() => ({}), {
  setNavShow
})(Navbar);
