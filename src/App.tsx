import './global.custom.scss';

import React, { lazy } from 'react';
import { Route, Routes } from 'react-router';

import RequireAuth from '@/components/RequireAuth';

import WithLoading from './components/WithLoading';

const Login = lazy(
  () => import(/* webpackChunkName:'Login', webpackPrefetch:true */ '@/pages/Login')
);
const Admin = lazy(
  () => import(/* webpackChunkName:'Admin', webpackPrefetch:true */ '@/pages/Admin')
);

const App: React.FC = () => {
  return (
    <WithLoading>
      <Routes>
        <Route
          path='/login'
          element={
            <RequireAuth requireLogin={false} to='/home'>
              <Login />
            </RequireAuth>
          }
        />
        <Route
          path='/*'
          element={
            <RequireAuth requireLogin={true} to='/login'>
              <Admin />
            </RequireAuth>
          }
        />
      </Routes>
    </WithLoading>
  );
};

export default App;
