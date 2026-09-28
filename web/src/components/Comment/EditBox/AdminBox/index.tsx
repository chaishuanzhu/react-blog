import { useKeyPress, useSafeState } from 'ahooks';
import { message } from 'antd';
import classNames from 'classnames';
import React, { memo, useRef } from 'react';

import type { User } from '@/utils/api';
import { login, toApiError } from '@/utils/api';

import s from './index.scss';

interface Props {
  showAdmin?: boolean;
  setShowAdmin?: Function;
  onLogin?: (user: User) => void;
}

const AdminBox: React.FC<Props> = ({ showAdmin = false, setShowAdmin, onLogin }) => {
  const pwdRef = useRef(null);

  const [adminEmail, setAdminEmail] = useSafeState('');
  const [adminPwd, setAdminPwd] = useSafeState('');

  const hideAdmin = () => {
    setShowAdmin?.(false);
    setAdminEmail('');
    setAdminPwd('');
  };

  const adminLogin = async () => {
    try {
      const user = await login(adminEmail.trim(), adminPwd);
      message.success('登录成功！');
      onLogin?.(user);
      hideAdmin();
    } catch (err) {
      const { status } = toApiError(err);
      message.error(status === 429 ? '尝试次数过多，请稍后再试！' : '登录失败，请重试！');
    }
  };

  useKeyPress(13, adminLogin, {
    target: pwdRef
  });

  return (
    <div className={classNames(s.adminBox, { [s.showAdmin]: showAdmin })}>
      <div className={s.itemBox}>
        <div className={s.adminKey}>邮箱</div>
        <input
          type='text'
          className={s.adminValue}
          value={adminEmail}
          onChange={e => setAdminEmail(e.target.value)}
        />
      </div>
      <div className={s.itemBox}>
        <div className={s.adminKey}>密码</div>
        <input
          ref={pwdRef}
          type='password'
          className={s.adminValue}
          value={adminPwd}
          onChange={e => setAdminPwd(e.target.value)}
        />
      </div>
      <div className={classNames(s.itemBox, s.adminBtns)}>
        <div className={s.adminBtn} onClick={hideAdmin}>
          取消
        </div>
        <div className={s.adminBtn} onClick={adminLogin}>
          登录
        </div>
      </div>
    </div>
  );
};

export default memo(AdminBox);
