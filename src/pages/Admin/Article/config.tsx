import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm, Tag } from '@arco-design/web-react';
import dayjs from 'dayjs';
import React from 'react';

import TableTag from '@/components/TableTag';
import type { AdminArticle, NamedRef } from '@/utils/api';
import { blogLink, dateTimeFormat } from '@/utils/constant';

interface Props {
  isDraft: boolean;
  handleEdit: (id: number) => void;
  handleDelete: (id: number) => void;
}

export const useColumns = ({
  isDraft,
  handleEdit,
  handleDelete
}: Props): TableColumnProps<AdminArticle>[] => [
  {
    title: '标题',
    dataIndex: 'title',
    render: (title: string) => <strong>{title}</strong>
  },
  {
    title: isDraft ? '保存日期' : '发布日期',
    dataIndex: isDraft ? 'updatedAt' : 'publishedAt',
    render: (time: string) => <>{dayjs(time).format(dateTimeFormat)}</>
  },
  {
    title: '分类',
    dataIndex: 'category',
    render: (category: NamedRef | null) =>
      category ? <Tag color='#2db7f5'>{category.name}</Tag> : null
  },
  {
    title: '标签',
    dataIndex: 'tags',
    render: (tags: NamedRef[]) => <TableTag tags={tags.map(t => t.name)} />
  },
  {
    title: '操作',
    render: (_: unknown, { id }: AdminArticle) => (
      <>
        {!isDraft && (
          <Button
            type='primary'
            style={{ marginRight: 10 }}
            onClick={() => window.open(blogLink(`/post/${id}`))}
          >
            查看
          </Button>
        )}
        <Button type='primary' style={{ marginRight: 10 }} onClick={() => handleEdit(id)}>
          编辑
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该文章吗？'
          onOk={() => handleDelete(id)}
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
