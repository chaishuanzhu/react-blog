import React, { Suspense } from 'react';

import Loading from '@/components/Loading';

const WithLoading: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <Suspense fallback={<Loading />}>{children}</Suspense>;
};

export default WithLoading;
