import { Message } from '@arco-design/web-react';
import { useTitle } from 'ahooks';
import React, { useState } from 'react';

import type { ModalField } from '@/components/CustomModal';
import CustomModal from '@/components/CustomModal';
import MyTable from '@/components/MyTable';
import PageHeader from '@/components/PageHeader';
import type { FriendLink, FriendLinkInput } from '@/utils/api';
import { friendLinkApi } from '@/utils/api';
import { defaultPageSize, siteTitle } from '@/utils/constant';
import { mutate } from '@/utils/feedback';
import { useClientTable } from '@/utils/hooks/useClientTable';

import { Title } from '../titleConfig';
import { useColumns } from './config';

const emptyForm: FriendLinkInput = { name: '', url: '', avatar: '', description: '' };

const Link: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Link}`);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [id, setId] = useState(0);
  const [form, setForm] = useState(emptyForm);

  const { data, total, loading, refresh, page, setPage, handleDelete } = useClientTable(
    friendLinkApi.list,
    friendLinkApi.remove,
    defaultPageSize
  );

  const openModal = (item?: FriendLink) => {
    setId(item?.id ?? 0);
    setForm(
      item
        ? { name: item.name, url: item.url, avatar: item.avatar, description: item.description }
        : emptyForm
    );
    setIsModalOpen(true);
  };

  const field = (text: string, key: keyof FriendLinkInput, placeholder?: string): ModalField => ({
    text,
    placeholder,
    value: form[key],
    onChange: value => setForm(prev => ({ ...prev, [key]: value }))
  });

  const modalOk = async () => {
    if (!form.name.trim() || !form.url.trim()) {
      Message.info('请输入友链名称和链接！');
      return;
    }
    const ok = await mutate(
      () => (id ? friendLinkApi.update(id, form) : friendLinkApi.create(form)),
      id ? '修改成功！' : '添加成功！'
    );
    if (!ok) return;
    setIsModalOpen(false);
    refresh();
  };

  const columns = useColumns({ handleEdit: openModal, handleDelete });

  return (
    <>
      <PageHeader text='添加友链' onClick={() => openModal()} />
      <MyTable
        loading={loading}
        columns={columns}
        data={data}
        total={total}
        page={page}
        setPage={setPage}
      />
      <CustomModal
        isEdit={!!id}
        isModalOpen={isModalOpen}
        name='友链'
        modalOk={modalOk}
        modalCancel={() => setIsModalOpen(false)}
        fields={[
          field('名称', 'name'),
          field('链接', 'url', 'https://'),
          field('头像', 'avatar', '（可选）https://'),
          field('描述', 'description', '（可选）')
        ]}
      />
    </>
  );
};

export default Link;
