import { UserOutlined } from '@ant-design/icons';
import {
  useBoolean,
  useKeyPress,
  useLocalStorageState,
  useMount,
  useSafeState
} from 'ahooks';
import { message } from 'antd';
import classNames from 'classnames';
import React, { useRef } from 'react';
import { connect } from 'react-redux';

import { setAvatar, setEmail, setLink, setName } from '@/redux/actions';
import type { storeState } from '@/redux/interface';
import type { User } from '@/utils/api';
import { clearToken, getMe, getToken, postComment, toApiError } from '@/utils/api';

import AdminBox from './AdminBox';
import Emoji from './Emoji';
import s from './index.scss';
import PreShow from './PreShow';

interface Props {
  articleId?: number;
  onPosted?: () => void;
  isReply?: boolean;
  name?: string;
  link?: string;
  email?: string;
  avatar?: string;
  setAvatar?: Function;
  setEmail?: Function;
  setLink?: Function;
  setName?: Function;
  closeReply?: Function;
  className?: string;
  replyName?: string;
  parentId?: number;
}

const maxContentLength = 1000;
const reQQ = /^[1-9][0-9]{4,11}$/;
const reQQEmail = /^([1-9][0-9]{4,11})@qq\.com$/i;

const qqAvatar = (qq: string) => `https://q1.qlogo.cn/g?b=qq&nk=${qq}&s=100`;

const EditBox: React.FC<Props> = ({
  articleId,
  onPosted,
  isReply = false,
  name,
  link,
  email,
  avatar,
  setAvatar,
  setEmail,
  setLink,
  setName,
  closeReply,
  replyName,
  parentId,
  className
}) => {
  const nameRef = useRef(null);

  const [showAdmin, setShowAdmin] = useSafeState(false);
  const [showPre, { toggle: togglePre, setFalse: closePre }] = useBoolean(false);
  const [sending, setSending] = useSafeState(false);

  const [text, setText] = useSafeState('');

  const [localName, setLocalName] = useLocalStorageState<string>('name');
  const [localEmail, setLocalEmail] = useLocalStorageState<string>('email');
  const [localLink, setLocalLink] = useLocalStorageState<string>('link');

  const isAdmin = !!getToken();
  const qqMatch = reQQEmail.exec(email || '');
  const previewAvatar = isAdmin ? avatar : qqMatch ? qqAvatar(qqMatch[1]) : '';

  const applyAdmin = (user: User) => {
    setName?.(user.nickname);
    setEmail?.(user.email);
    setLink?.(user.website);
    setAvatar?.(user.avatar);
  };

  const applyVisitor = () => {
    setName?.(localName || '');
    setEmail?.(localEmail || '');
    setLink?.(localLink || '');
    setAvatar?.('');
  };

  const validateConfig = [
    {
      check: /^[\u4e00-\u9fa5A-Za-z0-9_\- ]{2,16}$/,
      content: name?.trim(),
      errText: '昵称仅限中文、数字、字母、下划线、短横线，长度2~16！'
    },
    {
      check: /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/,
      content: email?.trim(),
      errText: '请输入正确的邮箱地址！'
    },
    {
      check: /^$|^https?:\/\/\S+$/,
      content: link?.trim(),
      errText: '网址需以 http:// 或 https:// 开头，或不填！'
    }
  ];

  const validate = () => {
    if (!text.trim()) {
      message.error('请输入内容再发布~');
      return false;
    }
    if (Array.from(text).length > maxContentLength) {
      message.error(`内容不能超过${maxContentLength}字！`);
      return false;
    }
    if (isAdmin) return true;
    const failed = validateConfig.find(({ check, content }) => !check.test(content || ''));
    if (failed) {
      message.error(failed.errText);
      return false;
    }
    return true;
  };

  const submit = async () => {
    if (sending || !validate()) return;
    setSending(true);
    try {
      await postComment({
        articleId,
        parentId,
        nickname: name?.trim() || '',
        email: email?.trim() || '',
        website: link?.trim() || '',
        content: text
      });
      setText('');
      closePre();
      closeReply?.();
      onPosted?.();
      message.success(`${isReply ? '回复' : '发布'}${articleId ? '评论' : '留言'}成功！`);
    } catch (err) {
      const { status, message: errMsg } = toApiError(err);
      if (status === 401) {
        message.warning('登录已过期，请重新登录！');
        applyVisitor();
      } else if (status === 403) {
        message.warning('不能使用站长的昵称或邮箱哦~');
      } else if (status === 429) {
        message.warning('发送太频繁了，请稍后再试~');
      } else {
        message.error(`发布失败：${errMsg}`);
      }
    } finally {
      setSending(false);
    }
  };

  useMount(() => {
    if (isReply) return;
    if (isAdmin) {
      getMe().then(applyAdmin, () => {
        clearToken();
        applyVisitor();
      });
      return;
    }
    applyVisitor();
  });

  const handleName = () => {
    if (isAdmin) return;
    if (name === 'admin') {
      setShowAdmin(true);
      setName?.('');
      return;
    }
    if (reQQ.test(name || '')) {
      const QQEmail = `${name}@qq.com`;
      setEmail?.(QQEmail);
      setLocalEmail(QQEmail);
      setName?.('');
      return;
    }
    setLocalName(name || '');
  };

  const logout = () => {
    clearToken();
    applyVisitor();
    message.success('已退出登录');
  };

  useKeyPress(13, handleName, {
    target: nameRef
  });

  const openPreShow = () => {
    if (!showPre && !text) {
      message.info('请写点什么再预览~');
      return;
    }
    togglePre();
  };

  const handleCloseReply = () => {
    closeReply?.();
  };

  return (
    <div className={classNames(s.editBox, className)}>
      {isReply && (
        <div className={s.replyNameBox}>
          回复给「<span>{replyName}</span>」：
        </div>
      )}
      <div className={s.flex}>
        <AdminBox showAdmin={showAdmin} setShowAdmin={setShowAdmin} onLogin={applyAdmin} />

        <div className={s.avatarBoxCol}>
          <div
            className={s.avatarBox}
            title={isAdmin ? '点击退出登录' : undefined}
            onClick={isAdmin ? logout : undefined}
          >
            {previewAvatar ? (
              <img src={previewAvatar} className={s.editAvatar} />
            ) : (
              <UserOutlined className={s.noAvatar} />
            )}
          </div>
        </div>
        <div className={s.editInputBox}>
          <div className={s.inputBox}>
            <div className={classNames(s.inputInfo, s.flex2)}>
              <div className={s.inputKey}>昵称</div>
              <input
                ref={nameRef}
                type='text'
                className={s.inputValue}
                placeholder='QQ号'
                value={name}
                readOnly={isAdmin}
                onChange={e => setName?.(e.target.value)}
                onBlur={handleName}
              />
            </div>
            <div className={classNames(s.inputInfo, s.flex3)}>
              <div className={s.inputKey}>邮箱</div>
              <input
                type='text'
                className={s.inputValue}
                placeholder='必填'
                value={email}
                readOnly={isAdmin}
                onChange={e => setEmail?.(e.target.value)}
                onBlur={e => !isAdmin && setLocalEmail(e.target.value)}
              />
            </div>
            <div className={classNames(s.inputInfo, s.flex3)}>
              <div className={s.inputKey}>网址</div>
              <input
                type='text'
                className={s.inputValue}
                placeholder='选填'
                value={link}
                readOnly={isAdmin}
                onChange={e => setLink?.(e.target.value)}
                onBlur={e => !isAdmin && setLocalLink(e.target.value)}
              />
            </div>
          </div>

          <div className={s.textareaBox}>
            <textarea
              className={s.textarea}
              value={text}
              onChange={e => setText(e.target.value)}
              placeholder='写点什么吗？支持markdown格式！&#10;可以在「昵称」处填写QQ号，自动获取「头像」和「QQ邮箱」！'
            />
          </div>
          <div className={s.commentBtns}>
            <Emoji />
            {isReply && (
              <div className={s.cancelBtn} onClick={handleCloseReply}>
                取消
              </div>
            )}
            <div className={s.previewBtn} onClick={openPreShow}>
              预览
            </div>
            <div className={s.sendBtn} onClick={submit}>
              {sending ? '发送中' : isReply ? '回复' : ' 发布'}
            </div>
          </div>
        </div>
      </div>
      <PreShow
        closePre={closePre}
        content={text}
        className={classNames({ [s.preShowHidden]: !showPre })}
      />
    </div>
  );
};

export default connect(
  (state: storeState) => ({
    name: state.name,
    link: state.link,
    email: state.email,
    avatar: state.avatar
  }),
  {
    setAvatar,
    setEmail,
    setLink,
    setName
  }
)(EditBox);
