import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm } from '@arco-design/web-react';
import React from 'react';

import type { FriendLink } from '@/utils/api';

interface Props {
  handleEdit: (item: FriendLink) => void;
  handleDelete: (id: number) => void;
}

export const useColumns = ({ handleEdit, handleDelete }: Props): TableColumnProps<FriendLink>[] => [
  {
    title: '名称',
    dataIndex: 'name'
  },
  {
    title: '链接',
    dataIndex: 'url',
    render: (text: string) => (
      <a href={text} target='_blank' rel='noreferrer'>
        {text}
      </a>
    )
  },
  {
    title: '头像',
    dataIndex: 'avatar',
    render: (url: string) =>
      url ? <img src={url} alt='avatar' style={{ width: 40, height: 40, borderRadius: 4 }} /> : null
  },
  {
    title: '描述',
    dataIndex: 'description'
  },
  {
    title: '操作',
    render: (_: unknown, item: FriendLink) => (
      <>
        <Button type='primary' style={{ marginRight: 10 }} onClick={() => handleEdit(item)}>
          更新
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该友链吗？'
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
