import { useRequest } from 'ahooks';

import { mutate, pageAfterDelete } from '../feedback';
import { usePage } from './usePage';

// 接口一次返回全部数据，由前端分页
export const useClientTable = <T>(
  list: () => Promise<T[]>,
  remove: (id: number) => Promise<void>,
  pageSize: number
) => {
  const { page, setPage } = usePage();
  const { data: all = [], loading, refresh } = useRequest(list);

  const handleDelete = async (id: number) => {
    if (!(await mutate(() => remove(id), '删除成功！'))) return;
    setPage(pageAfterDelete(all.length, page, pageSize));
    refresh();
  };

  return {
    data: all.slice(pageSize * (page - 1), pageSize * page),
    total: all.length,
    loading,
    refresh,
    page,
    setPage,
    handleDelete
  };
};
