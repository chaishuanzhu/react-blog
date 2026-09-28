import type { TableColumnProps } from '@arco-design/web-react';
import { Button, Popconfirm, Popover } from '@arco-design/web-react';
import dayjs from 'dayjs';
import React from 'react';
import { IoImage } from 'react-icons/io5';

import type { Moment } from '@/utils/api';
import { dateTimeFormat } from '@/utils/constant';

import s from './index.scss';

interface Props {
  handleEdit: (item: Moment) => void;
  handleDelete: (id: number) => void;
  onClickImg: (url: string) => void;
}

export const useColumns = ({
  handleEdit,
  handleDelete,
  onClickImg
}: Props): TableColumnProps<Moment>[] => [
  {
    title: '发布日期',
    dataIndex: 'createdAt',
    render: (time: string) => <>{dayjs(time).format(dateTimeFormat)}</>
  },
  {
    title: '图片',
    dataIndex: 'images',
    render: (images: string[]) =>
      images.length ? (
        <Popover
          position='right'
          className={s.imgsPopover}
          content={
            <div className={s.imgsBox}>
              {images.map((url, index) => (
                <div key={index} className={s.imgDiv} onClick={() => onClickImg(url)}>
                  <img src={url} alt='img' className={s.img} />
                </div>
              ))}
            </div>
          }
          trigger='hover'
        >
          <div className={s.imgHover}>
            <IoImage />
          </div>
        </Popover>
      ) : null
  },
  {
    title: '说说内容',
    dataIndex: 'content',
    render: (content: string) => (
      <div style={{ width: '100%', height: '100%' }}>
        <div style={{ margin: 'auto', width: 500, whiteSpace: 'pre-wrap' }}>{content}</div>
      </div>
    )
  },
  {
    title: '操作',
    render: (_: unknown, item: Moment) => (
      <>
        <Button type='primary' style={{ marginRight: 10 }} onClick={() => handleEdit(item)}>
          修改
        </Button>
        <Popconfirm
          position='br'
          title='确定要删除该说说吗？'
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
