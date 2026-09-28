import 'dayjs/locale/zh-cn';

import { MessageOutlined } from '@ant-design/icons';
import { useBoolean } from 'ahooks';
import classNames from 'classnames';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';

import MarkDown from '@/components/MarkDown';
import { siteConfig } from '@/site.config';
import type { Comment } from '@/utils/api';
import { useLazyImg } from '@/utils/hooks/useLazyImg';

import EditBox from '../../EditBox';
import s from './index.scss';

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);

interface Props {
  comment: Comment;
  isReply?: boolean;
  articleId?: number;
  onPosted?: () => void;
}

const MsgItem: React.FC<Props> = ({ comment, isReply, articleId, onPosted }) => {
  const { id, avatar, website, nickname, createdAt, content, isAdmin } = comment;
  const [showReply, { toggle: toggleReply, setFalse: closeReply }] = useBoolean(false);
  const { imgRef, imgUrl } = useLazyImg(avatar, siteConfig.loadingImage);

  return (
    <div
      className={classNames(s.commentItem, {
        [s.marginLeft]: isReply
      })}
    >
      <div className={s.flex}>
        <div ref={imgRef} className={s.avatarBox}>
          <img
            src={imgUrl}
            className={classNames({
              [s.avatar]: imgUrl !== siteConfig.loadingImage,
              [s.loading]: imgUrl === siteConfig.loadingImage
            })}
          />
        </div>
        {!isReply && (
          <div className={s.replyBtn} onClick={toggleReply}>
            <MessageOutlined />
          </div>
        )}

        <div className={s.contentBox}>
          <div className={s.usrInfo}>
            <a
              href={website || undefined}
              target={website ? '_blank' : '_self'}
              rel='noreferrer nofollow'
              className={s.name}
              style={{ cursor: website ? 'pointer' : 'default' }}
            >
              {nickname}
            </a>
            {isAdmin && <span className={s.flag}>站长</span>}
            <span className={s.date}>{dayjs(createdAt).fromNow()}</span>
          </div>
          <MarkDown content={content} className={s.content} />
        </div>
      </div>

      {!isReply && (
        <EditBox
          closeReply={closeReply}
          isReply={true}
          className={classNames(s.replyBox, { [s.replyHidden]: !showReply })}
          replyName={nickname}
          parentId={id}
          articleId={articleId}
          onPosted={onPosted}
        />
      )}
    </div>
  );
};

export default MsgItem;
