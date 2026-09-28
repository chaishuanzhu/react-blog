import { siteConfig } from '@/site.config';

// 路由前缀，如 '/admin'（来自 scripts/constant.js 的 PUBLIC_PATH）
export const basePath = (process.env.PUBLIC_PATH as string).replace(/\/$/, '');

// 标题、头像、博客地址在 src/site.config.ts 中修改
export const avatarUrl = siteConfig.defaultAvatar;
export const blogUrl = siteConfig.blogUrl;
export const blogLink = (path: string) => blogUrl.replace(/\/$/, '') + path;
export const siteTitle = siteConfig.title;

// 分页：默认每页数量
export const defaultPageSize = 12;
// 建站日志 分页
export const logPageSize = 10;

// 说说最多图片数（与服务端一致）
export const maxMomentImages = 9;

// 时间输入格式
export const dateTimeFormat = 'YYYY-MM-DD HH:mm:ss';
export const dateFormat = 'YYYY-MM-DD';
