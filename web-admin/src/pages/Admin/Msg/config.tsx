import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm } from '@arco-design/web-react';
import classNames from 'classnames';
import dayjs from 'dayjs';
import React from 'react';

import type { AdminComment } from '@/utils/api';
import { blogLink, dateTimeFormat } from '@/utils/constant';

import s from './index.scss';

interface Props {
  handleDelete: (id: number) => void;
}

export const useColumns = ({ handleDelete }: Props): TableColumnProps<AdminComment>[] => [
  {
    title: '昵称',
    dataIndex: 'nickname',
    render: (text: string, { isAdmin }: AdminComment) => (
      <div className={s.msgUserNameBox}>
        <div className={classNames(s.msgUserName, { [s.msgUserAdmin]: isAdmin })}>{text}</div>
      </div>
    )
  },
  {
    title: '联系邮箱',
    dataIndex: 'email'
  },
  {
    title: '网址',
    dataIndex: 'website',
    render: (text: string) => (
      <a href={text} target='_blank' rel='noreferrer'>
        {text}
      </a>
    )
  },
  {
    title: '日期',
    dataIndex: 'createdAt',
    render: (text: string) => <>{dayjs(text).format(dateTimeFormat)}</>
  },
  {
    title: '类型',
    render: (_: unknown, { article, parentId }: AdminComment) => (
      <div className={s.typeBox}>
        <div
          className={article ? s.comment : s.msg}
          style={parentId ? { marginRight: 5 } : {}}
          title={article?.title}
        >
          {article ? '文章评论' : '留言板'}
        </div>
        {parentId && <div className={s.reply}>回复</div>}
      </div>
    )
  },
  {
    title: '内容',
    dataIndex: 'content',
    width: 400,
    render: (text: string) => <div className={s.msgsContent}>{text}</div>
  },
  {
    title: 'IP',
    dataIndex: 'ip'
  },
  {
    title: '操作',
    render: (_: unknown, { article, id }: AdminComment) => (
      <>
        <Button
          style={{ marginRight: 10 }}
          type='primary'
          onClick={() => window.open(blogLink(article ? `/post/${article.id}` : '/msg'))}
        >
          查看
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该留言吗？其下回复会一并删除。'
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
