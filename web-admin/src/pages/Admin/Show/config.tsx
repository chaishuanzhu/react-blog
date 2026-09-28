import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm } from '@arco-design/web-react';
import React from 'react';

import type { Project } from '@/utils/api';

import s from './index.scss';

interface Props {
  handleEdit: (item: Project) => void;
  handleDelete: (id: number) => void;
  onClickImg: (url: string) => void;
}

export const useColumns = ({
  handleEdit,
  handleDelete,
  onClickImg
}: Props): TableColumnProps<Project>[] => [
  {
    title: '序号',
    dataIndex: 'sortOrder'
  },
  {
    title: '封面',
    dataIndex: 'cover',
    render: (url: string) =>
      url ? (
        <div className={s.tableCoverBox}>
          <img src={url} alt='cover' className={s.tableCover} onClick={() => onClickImg(url)} />
        </div>
      ) : null
  },
  {
    title: '名称',
    dataIndex: 'name'
  },
  {
    title: '描述',
    dataIndex: 'description'
  },
  {
    title: '操作',
    render: (_: unknown, item: Project) => (
      <>
        {item.url && (
          <Button
            type='primary'
            style={{ marginRight: 10 }}
            onClick={() => window.open(item.url)}
          >
            查看
          </Button>
        )}
        <Button style={{ marginRight: 10 }} type='primary' onClick={() => handleEdit(item)}>
          更新
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该作品吗？'
          onOk={() => handleDelete(item.id)}
          okText='Yes'
          cancelText='No'
        >
          <Button type='primary' status='danger'>
            删除
          </Button>
        </Popconfirm>
      </>
    )
  }
];
