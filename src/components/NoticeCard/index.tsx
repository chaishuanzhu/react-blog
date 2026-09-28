import { Input, Message } from '@arco-design/web-react';
import { IconLoading } from '@arco-design/web-react/icon';
import { useRequest } from 'ahooks';
import classNames from 'classnames';
import React, { useState } from 'react';

import { noticeApi } from '@/utils/api';
import { mutate } from '@/utils/feedback';

import CustomModal from '../CustomModal';
import Emoji from '../Emoji';
import s from './index.scss';

const { TextArea } = Input;

const NoticeCard: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [notice, setLocalNotice] = useState('');

  const { data, loading, refresh } = useRequest(noticeApi.get);

  const openModal = () => {
    setLocalNotice(data ?? '');
    setIsModalOpen(true);
  };

  const modalOk = async () => {
    if (!notice.trim()) {
      Message.warning('请输入公告内容~');
      return;
    }
    if (await mutate(() => noticeApi.update(notice.trim()), '修改成功！')) {
      setIsModalOpen(false);
      refresh();
    }
  };

  return (
    <>
      <div className={s.cardBox}>
        <div className={s.title}>公告</div>
        <div
          className={classNames(s.noticeText, { [s.loading]: loading })}
          onClick={openModal}
        >
          {loading ? <IconLoading /> : data}
        </div>
      </div>
      <CustomModal
        isEdit={true}
        isModalOpen={isModalOpen}
        name='公告'
        modalOk={modalOk}
        modalCancel={() => setIsModalOpen(false)}
      >
        <TextArea
          placeholder='请输入公告内容'
          maxLength={21 * 4}
          allowClear
          showWordLimit
          value={notice}
          onChange={value => setLocalNotice(value)}
          autoSize={false}
          style={{ height: 100, resize: 'none' }}
        />
        <Emoji style={{ marginTop: 10 }} />
      </CustomModal>
    </>
  );
};

export default NoticeCard;
