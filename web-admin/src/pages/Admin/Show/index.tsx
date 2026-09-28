import { Input, Message } from '@arco-design/web-react';
import { useTitle } from 'ahooks';
import React, { useState } from 'react';

import CustomModal from '@/components/CustomModal';
import ImgView from '@/components/ImgView';
import MyTable from '@/components/MyTable';
import PageHeader from '@/components/PageHeader';
import UploadButton from '@/components/UploadButton';
import type { Project } from '@/utils/api';
import { projectApi } from '@/utils/api';
import { showPageSize, siteTitle } from '@/utils/constant';
import { mutate } from '@/utils/feedback';
import { useClientTable } from '@/utils/hooks/useClientTable';

import { Title } from '../titleConfig';
import { useColumns } from './config';

interface Form {
  sortOrder: string;
  name: string;
  description: string;
  cover: string;
  url: string;
}

const emptyForm: Form = { sortOrder: '0', name: '', description: '', cover: '', url: '' };

const Show: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Show}`);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [id, setId] = useState(0);
  const [form, setForm] = useState(emptyForm);

  const [imgUrl, setImgUrl] = useState('');
  const [isViewShow, setIsViewShow] = useState(false);

  const { data, total, loading, refresh, page, setPage, handleDelete } = useClientTable(
    projectApi.list,
    projectApi.remove,
    showPageSize
  );

  const openModal = (item?: Project) => {
    setId(item?.id ?? 0);
    setForm(
      item
        ? {
            sortOrder: String(item.sortOrder),
            name: item.name,
            description: item.description,
            cover: item.cover,
            url: item.url
          }
        : emptyForm
    );
    setIsModalOpen(true);
  };

  const input = (text: string, key: keyof Form, placeholder?: string) => (
    <Input
      size='large'
      addBefore={text}
      placeholder={placeholder}
      value={form[key]}
      onChange={value => setForm(prev => ({ ...prev, [key]: value }))}
      style={{ marginBottom: 10 }}
    />
  );

  const modalOk = async () => {
    const sortOrder = Number(form.sortOrder);
    if (!form.name.trim() || !Number.isInteger(sortOrder)) {
      Message.info('请输入作品名称和整数序号！');
      return;
    }
    const body = { ...form, sortOrder };
    const ok = await mutate(
      () => (id ? projectApi.update(id, body) : projectApi.create(body)),
      id ? '修改成功！' : '添加成功！'
    );
    if (!ok) return;
    setIsModalOpen(false);
    refresh();
  };

  const columns = useColumns({
    handleEdit: openModal,
    handleDelete,
    onClickImg: (url: string) => {
      setIsViewShow(true);
      setImgUrl(url);
    }
  });

  return (
    <>
      <PageHeader text='添加作品' onClick={() => openModal()} />
      <MyTable
        loading={loading}
        columns={columns}
        data={data}
        total={total}
        page={page}
        pageSize={showPageSize}
        setPage={setPage}
      />
      <CustomModal
        isEdit={!!id}
        isModalOpen={isModalOpen}
        name='作品'
        modalOk={modalOk}
        modalCancel={() => setIsModalOpen(false)}
      >
        {input('序号', 'sortOrder', '越小越靠前')}
        {input('名称', 'name')}
        {input('描述', 'description', '（可选）')}
        <div style={{ display: 'flex', marginBottom: 10 }}>
          <Input
            size='large'
            addBefore='封面'
            placeholder='（可选）https://'
            value={form.cover}
            onChange={cover => setForm(prev => ({ ...prev, cover }))}
            style={{ marginRight: 10 }}
          />
          <UploadButton
            size='large'
            text='上传'
            onUploaded={cover => setForm(prev => ({ ...prev, cover }))}
          />
        </div>
        {input('链接', 'url', '（可选）https://')}
      </CustomModal>
      <ImgView
        isViewShow={isViewShow}
        viewUrl={imgUrl}
        onClick={() => setIsViewShow(false)}
      />
    </>
  );
};

export default Show;
