import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm } from '@arco-design/web-react';
import dayjs from 'dayjs';
import React from 'react';

import type { Changelog } from '@/utils/api';
import { dateFormat } from '@/utils/constant';

interface Props {
  handleEdit: (item: Changelog) => void;
  handleDelete: (id: number) => void;
}

export const useColumns = ({ handleEdit, handleDelete }: Props): TableColumnProps<Changelog>[] => [
  {
    title: '日期',
    dataIndex: 'loggedAt',
    render: (time: string) => <>{dayjs(time).format(dateFormat)}</>
  },
  {
    title: '日志内容',
    dataIndex: 'items',
    render: (items: string[]) => items.map((item, index) => <div key={index}>{item}</div>)
  },
  {
    title: '操作',
    render: (_: unknown, item: Changelog) => (
      <>
        <Button style={{ marginRight: 10 }} type='primary' onClick={() => handleEdit(item)}>
          修改
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该日志吗？'
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
