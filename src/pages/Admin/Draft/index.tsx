import { useTitle } from 'ahooks';
import React from 'react';

import MyTable from '@/components/MyTable';
import { siteTitle } from '@/utils/constant';

import { useArticleTable } from '../Article/useArticleTable';
import { Title } from '../titleConfig';

const Draft: React.FC = () => {
  useTitle(`${siteTitle} | ${Title.Drafts}`);

  const { columns, data, total, loading, page, setPage } = useArticleTable('draft');

  return (
    <MyTable
      loading={loading}
      columns={columns}
      data={data}
      total={total}
      page={page}
      setPage={setPage}
      noHeader={true}
    />
  );
};

export default Draft;
