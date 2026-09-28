import { useRequest } from 'ahooks';
import { useNavigate } from 'react-router';

import type { ArticleQuery, ArticleStatus } from '@/utils/api';
import { articleApi } from '@/utils/api';
import { defaultPageSize } from '@/utils/constant';
import { mutate, pageAfterDelete } from '@/utils/feedback';
import { usePage } from '@/utils/hooks/usePage';

import { useColumns } from './config';

type Filter = Omit<ArticleQuery, 'status' | 'page' | 'pageSize'>;

export const useArticleTable = (status: ArticleStatus, filter: Filter = {}) => {
  const navigate = useNavigate();
  const { page, setPage } = usePage();

  const { data, loading, refresh } = useRequest(
    () => articleApi.list({ status, page, pageSize: defaultPageSize, ...filter }),
    { refreshDeps: [status, page, filter.keyword, filter.categoryId, filter.tagId] }
  );

  const handleDelete = async (id: number) => {
    if (!(await mutate(() => articleApi.remove(id), '删除成功！'))) return;
    const next = pageAfterDelete(data?.total ?? 0, page, defaultPageSize);
    if (next === page) refresh();
    else setPage(next);
  };

  const columns = useColumns({
    isDraft: status === 'draft',
    handleEdit: id => navigate(`/addArticle?id=${id}`),
    handleDelete
  });

  return {
    columns,
    data: data?.items ?? [],
    total: data?.total ?? 0,
    loading,
    page,
    setPage
  };
};
