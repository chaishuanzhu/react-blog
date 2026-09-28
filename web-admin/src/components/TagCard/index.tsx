import { Input, Message, Popconfirm } from '@arco-design/web-react';
import { IconDelete, IconEdit, IconLoading } from '@arco-design/web-react/icon';
import { useRequest } from 'ahooks';
import classNames from 'classnames';
import React, { useState } from 'react';
import { useNavigate } from 'react-router';

import { tagApi } from '@/utils/api';
import { mutate } from '@/utils/feedback';

import CustomModal from '../CustomModal';
import { useColor } from './config';
import s from './index.scss';

const { Search } = Input;

const TagCard: React.FC = () => {
  const navigate = useNavigate();
  const [editing, setEditing] = useState<{ id: number; name: string } | null>(null);
  const [newTag, setNewTag] = useState('');

  const { data: tags = [], loading, refresh } = useRequest(tagApi.list);
  const { tagColor, colorLen } = useColor();

  const isExist = (name: string, exceptId = 0) =>
    tags.some(t => t.name === name && t.id !== exceptId);

  const modalOk = async () => {
    if (!editing) return;
    const name = editing.name.trim();
    if (!name) {
      Message.warning('请输入标签名称~');
      return;
    }
    if (isExist(name, editing.id)) {
      Message.warning('标签名称已存在~');
      return;
    }
    if (await mutate(() => tagApi.update(editing.id, { name }), '修改成功！')) {
      setEditing(null);
      refresh();
    }
  };

  const addNewTag = async () => {
    const name = newTag.trim();
    if (!name) {
      Message.warning('请输入标签名称~');
      return;
    }
    if (isExist(name)) {
      Message.warning('标签名称已存在~');
      return;
    }
    if (await mutate(() => tagApi.create({ name }), '添加成功！')) {
      setNewTag('');
      refresh();
    }
  };

  const deleteTag = async (id: number) => {
    if (await mutate(() => tagApi.remove(id), '删除成功！')) {
      refresh();
    }
  };

  return (
    <>
      <div className={s.cardBox}>
        <div className={s.title}>标签</div>
        <Search
          size='default'
          allowClear
          placeholder='新建标签'
          searchButton='创建'
          value={newTag}
          onChange={(value: string) => setNewTag(value)}
          onSearch={addNewTag}
        />
        <div className={classNames(s.tagsBox, { [s.tagLoading]: loading })}>
          {loading ? (
            <IconLoading />
          ) : (
            tags.map(({ id, name, articleCount }, index) => (
              <div
                key={id}
                className={s.tagItem}
                style={{ backgroundColor: tagColor[index % colorLen] }}
                title={`${articleCount} 篇文章，双击查看`}
                onDoubleClick={() => navigate(`/article?tagId=${id}`)}
              >
                {name}
                <IconEdit className={s.iconBtn} onClick={() => setEditing({ id, name })} />
                <Popconfirm
                  position='br'
                  title={`确定要删除「${name}」吗？`}
                  onOk={() => deleteTag(id)}
                  okText='Yes'
                  cancelText='No'
                >
                  <IconDelete className={s.iconBtn} />
                </Popconfirm>
              </div>
            ))
          )}
        </div>
      </div>
      <CustomModal
        isEdit={true}
        isModalOpen={!!editing}
        name='标签'
        modalOk={modalOk}
        modalCancel={() => setEditing(null)}
      >
        <Input
          size='default'
          value={editing?.name ?? ''}
          onChange={name => setEditing(prev => prev && { ...prev, name })}
          onPressEnter={modalOk}
        />
      </CustomModal>
    </>
  );
};

export default TagCard;
