import { Input, Message } from '@arco-design/web-react';
import { useTitle } from 'ahooks';
import dayjs from 'dayjs';
import React, { useState } from 'react';

import CustomModal from '@/components/CustomModal';
import Emoji from '@/components/Emoji';
import MyTable from '@/components/MyTable';
import PageHeader from '@/components/PageHeader';
import type { Changelog } from '@/utils/api';
import { changelogApi } from '@/utils/api';
import { dateFormat, logPageSize, siteTitle } from '@/utils/constant';
import { mutate, parseLocalTime } from '@/utils/feedback';
import { useClientTable } from '@/utils/hooks/useClientTable';

import { Title } from '../titleConfig';
import { useColumns } from './config';

const { TextArea } = Input;

const Log: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Log}`);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [id, setId] = useState(0);
  const [date, setDate] = useState('');
  const [text, setText] = useState('');

  const { data, total, loading, refresh, page, setPage, handleDelete } = useClientTable(
    changelogApi.list,
    changelogApi.remove,
    logPageSize
  );

  const openModal = (item?: Changelog) => {
    setId(item?.id ?? 0);
    setDate(dayjs(item?.loggedAt).format(dateFormat));
    setText(item?.items.join('\n') ?? '');
    setIsModalOpen(true);
  };

  const modalOk = async () => {
    const loggedAt = parseLocalTime(date, dateFormat);
    const items = text
      .split('\n')
      .map(line => line.trim())
      .filter(Boolean);
    if (!loggedAt || !items.length) {
      Message.info('请输入合法的日期和日志内容！');
      return;
    }
    const ok = await mutate(
      () =>
        id ? changelogApi.update(id, { items, loggedAt }) : changelogApi.create({ items, loggedAt }),
      id ? '修改成功！' : '添加成功！'
    );
    if (!ok) return;
    setIsModalOpen(false);
    refresh();
  };

  const columns = useColumns({ handleEdit: openModal, handleDelete });

  return (
    <>
      <PageHeader text='添加日志' onClick={() => openModal()} />
      <MyTable
        loading={loading}
        columns={columns}
        data={data}
        total={total}
        page={page}
        pageSize={logPageSize}
        setPage={setPage}
      />
      <CustomModal
        isEdit={!!id}
        isModalOpen={isModalOpen}
        name='日志'
        modalOk={modalOk}
        modalCancel={() => setIsModalOpen(false)}
        updateText='修改'
      >
        <Input
          size='large'
          addBefore='时间'
          style={{ marginBottom: 10 }}
          value={date}
          onChange={value => setDate(value)}
        />
        <TextArea
          style={{ resize: 'none', marginBottom: 10, height: 120 }}
          placeholder='请输入日志内容，回车分隔'
          value={text}
          onChange={value => setText(value)}
        />
        <Emoji />
      </CustomModal>
    </>
  );
};

export default Log;
