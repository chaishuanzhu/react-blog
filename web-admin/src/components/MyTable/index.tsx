import type { TableColumnProps } from '@arco-design/web-react';
import { Pagination, Table } from '@arco-design/web-react';
import classNames from 'classnames';
import React from 'react';

import { defaultPageSize } from '@/utils/constant';

import s from './index.scss';

interface Props<T> {
  loading: boolean;
  columns: TableColumnProps<T>[];
  data: T[];
  total: number;
  page: number;
  pageSize?: number;
  noHeader?: boolean;
  setPage: (page: number) => void;
}

const MyTable = <T extends { id: number }>({
  loading,
  columns,
  data,
  total,
  page,
  pageSize = defaultPageSize,
  noHeader = false,
  setPage
}: Props<T>) => (
  <>
    <div className={classNames(s.myTableBox, { [s.noHeader]: noHeader })}>
      <Table
        border
        borderCell
        loading={loading}
        columns={columns}
        data={data}
        rowKey='id'
        showSorterTooltip={false}
        className={s.myTable}
        pagination={false}
      />
    </div>
    <div className={s.paginationBox}>
      <Pagination
        size='large'
        current={page}
        total={total}
        pageSize={pageSize}
        sizeCanChange={false}
        onChange={(page: number) => setPage(page)}
        hideOnSinglePage={true}
        showTotal={true}
      />
    </div>
  </>
);

export default MyTable;
