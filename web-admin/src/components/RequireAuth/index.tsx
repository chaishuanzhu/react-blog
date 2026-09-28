import React from 'react';
import { Navigate, useLocation } from 'react-router';

import { getToken } from '@/utils/api';

interface Props {
  requireLogin: boolean;
  to: string;
  children: React.ReactNode;
}

const RequireAuth: React.FC<Props> = ({ requireLogin, to, children }) => {
  const location = useLocation();
  const isLogin = !!getToken();

  if (requireLogin === isLogin) return children;
  return <Navigate to={to} state={{ from: location }} replace />;
};

export default RequireAuth;
