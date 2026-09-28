import classNames from 'classnames';
import React from 'react';
import { Navigate, NavLink, Route, Routes } from 'react-router';

import Header from '@/components/Header';
import WithLoading from '@/components/WithLoading';

import type { RouteType } from './config';
import { useRoutes } from './config';
import s from './index.scss';

const Admin: React.FC = () => {
  const routes = useRoutes();

  return (
    <div className={s.adminBox}>
      <Header />
      <div className={s.mainContent}>
        {/* 左侧导航区 */}
        <nav className={s.leftNav}>
          {routes.map((item: RouteType) => {
            return item.disPlayName ? (
              <div key={item.path} className={s.liItem}>
                <NavLink
                  to={`/${item.path}`}
                  className={({ isActive }) =>
                    isActive ? classNames(s.navItem, s.navItemActive) : s.navItem
                  }
                >
                  <div className={s.navIcon}>{item.icon}</div>
                  {item.disPlayName}
                </NavLink>
              </div>
            ) : null;
          })}
        </nav>
        {/* 右侧内容区域 */}
        <div className={s.rightContent}>
          <WithLoading>
            <Routes>
              {routes.map((item: RouteType) => (
                <Route key={item.path} path={item.path} element={item.element} />
              ))}
              <Route path='*' element={<Navigate to='/home' replace />} />
            </Routes>
          </WithLoading>
        </div>
      </div>
    </div>
  );
};

export default Admin;
