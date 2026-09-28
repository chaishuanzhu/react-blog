import React from 'react';

import LayoutLoading from '@/components/LayoutLoading';
import type { CommentThread } from '@/utils/api';

import s from './index.scss';
import MsgItem from './MsgItem';

interface Props {
  threads?: CommentThread[];
  loading?: boolean;
  articleId?: number;
  onPosted?: () => void;
}

const MsgList: React.FC<Props> = ({ threads, loading, articleId, onPosted }) => {
  return (
    <>
      {loading ? (
        <LayoutLoading />
      ) : (
        threads?.map(thread => (
          <div key={thread.id} className={s.completeMsg}>
            <MsgItem
              comment={thread}
              isReply={false}
              articleId={articleId}
              onPosted={onPosted}
            />
            {thread.replies.map(reply => (
              <MsgItem key={reply.id} comment={reply} isReply={true} />
            ))}
          </div>
        ))
      )}
    </>
  );
};

export default MsgList;
