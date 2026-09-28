import { useRequest, useSafeState } from 'ahooks';
import React from 'react';

import { getComments } from '@/utils/api';
import { msgSize } from '@/utils/constant';

import MyPagination from '../MyPagination';
import Divider from './Divider';
import EditBox from './EditBox';
import MsgList from './MsgList';
import Placehold from './Placehold';

interface Props {
  articleId?: number;
  autoScroll?: boolean;
  scrollToTop?: number;
}

const Comment: React.FC<Props> = ({ articleId, autoScroll = false, scrollToTop = 0 }) => {
  const [page, setPage] = useSafeState(1);

  const { data, loading, refresh } = useRequest(
    () => getComments(articleId, page, msgSize),
    {
      retryCount: 3,
      refreshDeps: [page, articleId]
    }
  );

  return (
    <div>
      <Divider />
      <EditBox articleId={articleId} onPosted={refresh} />
      <Placehold msgCount={data?.total} isMsg={!articleId} />
      <MsgList
        threads={data?.items}
        loading={loading && !data}
        articleId={articleId}
        onPosted={refresh}
      />
      <MyPagination
        current={page}
        defaultPageSize={msgSize}
        total={data?.total}
        setPage={setPage}
        autoScroll={autoScroll}
        scrollToTop={scrollToTop}
      />
    </div>
  );
};

export default Comment;
