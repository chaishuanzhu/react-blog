import { Input, Message } from '@arco-design/web-react';
import { useRequest, useTitle } from 'ahooks';
import dayjs from 'dayjs';
import React, { useState } from 'react';

import CustomModal from '@/components/CustomModal';
import Emoji from '@/components/Emoji';
import ImgView from '@/components/ImgView';
import MyTable from '@/components/MyTable';
import PageHeader from '@/components/PageHeader';
import UploadButton from '@/components/UploadButton';
import type { Moment } from '@/utils/api';
import { momentApi } from '@/utils/api';
import { dateTimeFormat, defaultPageSize, maxMomentImages, siteTitle } from '@/utils/constant';
import { mutate, pageAfterDelete, parseLocalTime } from '@/utils/feedback';
import { usePage } from '@/utils/hooks/usePage';

import { Title } from '../titleConfig';
import { useColumns } from './config';

const { TextArea } = Input;

const Say: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Say}`);

  const { page, setPage } = usePage();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [id, setId] = useState(0);
  const [date, setDate] = useState('');
  const [content, setContent] = useState('');
  const [images, setImages] = useState<string[]>([]);
  const [saving, setSaving] = useState(false);

  const [imgUrl, setImgUrl] = useState('');
  const [isViewShow, setIsViewShow] = useState(false);

  const { data, loading, refresh } = useRequest(
    () => momentApi.list(page, defaultPageSize),
    { refreshDeps: [page] }
  );

  const openModal = (item?: Moment) => {
    setId(item?.id ?? 0);
    setDate(dayjs(item?.createdAt).format(dateTimeFormat));
    setContent(item?.content ?? '');
    setImages(item?.images ?? []);
    setIsModalOpen(true);
  };

  const handleDelete = async (momentId: number) => {
    if (!(await mutate(() => momentApi.remove(momentId), '删除成功！'))) return;
    const next = pageAfterDelete(data?.total ?? 0, page, defaultPageSize);
    if (next === page) refresh();
    else setPage(next);
  };

  const columns = useColumns({
    handleEdit: openModal,
    handleDelete,
    onClickImg: (url: string) => {
      setIsViewShow(true);
      setImgUrl(url);
    }
  });

  const modalOk = async () => {
    const createdAt = parseLocalTime(date, dateTimeFormat);
    if (!createdAt || !content.trim()) {
      Message.info('请输入合法的时间和说说内容！');
      return;
    }
    const input = {
      content: content.trim(),
      images: images.map(url => url.trim()).filter(Boolean),
      createdAt
    };
    setSaving(true);
    const ok = await mutate(
      () => (id ? momentApi.update(id, input) : momentApi.create(input)),
      id ? '修改成功！' : '发表成功！'
    );
    setSaving(false);
    if (!ok) return;
    setIsModalOpen(false);
    if (id || page === 1) refresh();
    else setPage(1);
  };

  return (
    <>
      <PageHeader text='发表说说' onClick={() => openModal()} />
      <MyTable
        loading={loading}
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        page={page}
        setPage={setPage}
      />
      <CustomModal
        isEdit={!!id}
        isModalOpen={isModalOpen}
        name='说说'
        modalOk={modalOk}
        modalCancel={() => setIsModalOpen(false)}
        confirmLoading={saving}
        addText='发表'
        updateText='修改'
      >
        <Input
          size='large'
          addBefore='时间'
          value={date}
          style={{ marginBottom: 10 }}
          onChange={value => setDate(value)}
        />
        <TextArea
          placeholder='说说内容'
          style={{ resize: 'none', marginBottom: 10, height: 100 }}
          value={content}
          onChange={value => setContent(value)}
        />
        <TextArea
          style={{ resize: 'none', marginBottom: 10, height: 98 }}
          placeholder={`（可选）图片url，回车分隔，最多${maxMomentImages}张`}
          value={images.join('\n')}
          onChange={value => setImages(value.split('\n').slice(0, maxMomentImages))}
        />
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Emoji />
          <UploadButton
            style={{ marginLeft: 'auto' }}
            onUploaded={url =>
              setImages(prev =>
                [...prev.filter(Boolean), url].slice(0, maxMomentImages)
              )
            }
          />
        </div>
      </CustomModal>
      <ImgView
        isViewShow={isViewShow}
        viewUrl={imgUrl}
        onClick={() => setIsViewShow(false)}
      />
    </>
  );
};

export default Say;
