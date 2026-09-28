import { Dropdown, Menu, Message, Popconfirm } from '@arco-design/web-react';
import { useRequest } from 'ahooks';
import React, { useState } from 'react';
import { IoHome, IoLogOut, IoSettingsSharp } from 'react-icons/io5';
import { useNavigate } from 'react-router';

import { clearToken, getMe } from '@/utils/api';
import { avatarUrl, blogUrl } from '@/utils/constant';
import { useTime } from '@/utils/hooks/useTime';

import s from './index.scss';
import PasswordModal from './PasswordModal';
import ProfileModal from './ProfileModal';

type Dialog = 'profile' | 'password' | null;

const Header: React.FC = () => {
  const navigate = useNavigate();
  const { timeText } = useTime();
  const { data: me, mutate } = useRequest(getMe);
  const [dialog, setDialog] = useState<Dialog>(null);

  const logout = () => {
    clearToken();
    Message.success('已退出个人博客后台管理系统！');
    navigate('/login');
  };

  const menu = (
    <Menu onClickMenuItem={key => setDialog(key as Dialog)}>
      <Menu.Item key='profile'>个人资料</Menu.Item>
      <Menu.Item key='password'>修改密码</Menu.Item>
    </Menu>
  );

  return (
    <div className={s.headerBox}>
      <img src={me?.avatar || avatarUrl} alt='' className={s.avatar} />
      <div className={s.avatarText}>
        {timeText}，<span className={s.userName}>{me?.nickname ?? ''}</span>！
      </div>
      <a className={s.blogBtn} href={blogUrl} target='_blank' rel='noreferrer'>
        <IoHome />
      </a>
      <Dropdown droplist={menu} position='br' trigger='click'>
        <div className={s.settingBtn}>
          <IoSettingsSharp />
        </div>
      </Dropdown>
      <Popconfirm
        title='确定退出吗？'
        position='br'
        onOk={logout}
        okText='Yes'
        cancelText='No'
      >
        <div className={s.logoutBtn}>
          <IoLogOut />
        </div>
      </Popconfirm>
      <ProfileModal
        visible={dialog === 'profile'}
        user={me}
        onClose={() => setDialog(null)}
        onSaved={mutate}
      />
      <PasswordModal visible={dialog === 'password'} onClose={() => setDialog(null)} />
    </div>
  );
};

export default Header;
