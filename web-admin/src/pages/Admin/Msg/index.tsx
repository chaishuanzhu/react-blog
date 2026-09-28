import { useRequest, useTitle } from 'ahooks';
import React from 'react';

import MyTable from '@/components/MyTable';
import { commentApi } from '@/utils/api';
import { defaultPageSize, siteTitle } from '@/utils/constant';
import { mutate, pageAfterDelete } from '@/utils/feedback';
import { usePage } from '@/utils/hooks/usePage';

import { Title } from '../titleConfig';
import { useColumns } from './config';

const Msg: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Msg}`);
  const { page, setPage } = usePage();

  const { data, loading, refresh } = useRequest(
    () => commentApi.list(page, defaultPageSize),
    { refreshDeps: [page] }
  );

  const handleDelete = async (id: number) => {
    if (!(await mutate(() => commentApi.remove(id), '删除成功！'))) return;
    const next = pageAfterDelete(data?.total ?? 0, page, defaultPageSize);
    if (next === page) refresh();
    else setPage(next);
  };

  const columns = useColumns({ handleDelete });

  return (
    <MyTable
      loading={loading}
      columns={columns}
      data={data?.items ?? []}
      total={data?.total ?? 0}
      page={page}
      setPage={setPage}
      noHeader={true}
    />
  );
};

export default Msg;
