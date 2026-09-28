import { Button, Input, Message } from '@arco-design/web-react';
import { useTitle } from 'ahooks';
import React, { useState } from 'react';
import { BiLockAlt, BiUser } from 'react-icons/bi';
import { useNavigate } from 'react-router';

import { login } from '@/utils/api';
import { avatarUrl, siteTitle } from '@/utils/constant';
import { showError } from '@/utils/feedback';

import s from './index.scss';

const Login: React.FC = () => {
  useTitle(siteTitle);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  const handleLogin = async () => {
    if (!email || !password) {
      Message.warning('登录失败！请输入账号、密码！');
      return;
    }
    setLoading(true);
    try {
      await login(email.trim(), password);
      Message.success('登录成功！欢迎进入个人博客后台管理系统！');
      navigate('/home', { replace: true });
    } catch (err) {
      showError(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={s.LoginBox}>
      <div className={s.leftBox}>个人博客后台管理系统</div>
      <div className={s.rightBox}>
        <div className={s.avatarBox}>
          <img src={avatarUrl} alt='avatar' />
        </div>
        <div className={s.loginBox}>
          <Input
            style={{ marginBottom: 20 }}
            size='large'
            prefix={<BiUser />}
            placeholder='邮箱'
            autoComplete='username'
            value={email}
            onChange={value => setEmail(value)}
          />
          <Input.Password
            style={{ marginBottom: 20 }}
            size='large'
            prefix={<BiLockAlt />}
            placeholder='密码'
            autoComplete='current-password'
            value={password}
            onChange={value => setPassword(value)}
            onPressEnter={handleLogin}
          />
          <div className={s.btnBox}>
            <Button type='primary' size='large' long loading={loading} onClick={handleLogin}>
              登录
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
