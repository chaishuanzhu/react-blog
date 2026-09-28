import { Button, Input, Message, Popconfirm } from '@arco-design/web-react';
import { IconDelete, IconEdit, IconLoading } from '@arco-design/web-react/icon';
import classNames from 'classnames';
import React, { useState } from 'react';
import { useNavigate } from 'react-router';

import type { CategoryList } from '@/utils/api';
import { categoryApi } from '@/utils/api';
import { mutate } from '@/utils/feedback';

import CustomModal from '../CustomModal';
import s from './index.scss';

const { Search } = Input;

interface Props {
  categories?: CategoryList;
  loading: boolean;
  onChanged: () => void;
}

const ClassCard: React.FC<Props> = ({ categories, loading, onChanged }) => {
  const navigate = useNavigate();
  const [editing, setEditing] = useState<{ id: number; name: string } | null>(null);
  const [newName, setNewName] = useState('');

  const items = categories?.items ?? [];
  const isExist = (name: string, exceptId = 0) =>
    items.some(c => c.name === name && c.id !== exceptId);

  const modalOk = async () => {
    if (!editing) return;
    const name = editing.name.trim();
    if (!name) {
      Message.warning('请输入分类名称~');
      return;
    }
    if (isExist(name, editing.id)) {
      Message.warning('分类名称已存在~');
      return;
    }
    if (await mutate(() => categoryApi.update(editing.id, { name }), '修改成功！')) {
      setEditing(null);
      onChanged();
    }
  };

  const addNewClass = async () => {
    const name = newName.trim();
    if (!name) {
      Message.warning('请输入分类名称~');
      return;
    }
    if (isExist(name)) {
      Message.warning('分类名称已存在~');
      return;
    }
    if (await mutate(() => categoryApi.create({ name }), '添加成功！')) {
      setNewName('');
      onChanged();
    }
  };

  const deleteClass = async (id: number) => {
    if (await mutate(() => categoryApi.remove(id), '删除成功！')) {
      onChanged();
    }
  };

  const rows = [
    ...items,
    { id: 0, name: '未分类', articleCount: categories?.uncategorizedCount ?? 0 }
  ];

  return (
    <>
      <div className={s.cardBox}>
        <div className={s.title}>分类</div>
        <Search
          size='default'
          allowClear
          placeholder='新建分类'
          searchButton='创建'
          value={newName}
          onChange={(value: string) => setNewName(value)}
          onSearch={addNewClass}
        />
        <div className={classNames(s.classesBox, { [s.classLoading]: loading })}>
          {loading ? (
            <IconLoading />
          ) : (
            rows.map(({ id, name, articleCount }) => (
              <div key={id} className={s.classItem}>
                <div className={s.count}>{articleCount}</div>
                <div className={s.classTextBox}>
                  <div
                    className={s.classText}
                    onClick={() => id && navigate(`/article?categoryId=${id}`)}
                  >
                    《{name}》
                  </div>
                </div>
                <Button
                  type='primary'
                  className={s.classBtn}
                  icon={<IconEdit />}
                  onClick={() => setEditing({ id, name })}
                  disabled={!id}
                />
                <Popconfirm
                  position='br'
                  title={`确定要删除《${name}》吗？其下文章将变为未分类。`}
                  onOk={() => deleteClass(id)}
                  okText='Yes'
                  cancelText='No'
                  disabled={!id}
                >
                  <Button
                    style={{ width: 30, height: 30 }}
                    type='primary'
                    status='danger'
                    className={s.classBtn}
                    icon={<IconDelete />}
                    disabled={!id}
                  />
                </Popconfirm>
              </div>
            ))
          )}
        </div>
      </div>
      <CustomModal
        isEdit={true}
        isModalOpen={!!editing}
        name='分类'
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

export default ClassCard;
