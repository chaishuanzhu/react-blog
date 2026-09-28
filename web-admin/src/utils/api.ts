import type { AxiosError } from 'axios';
import axios from 'axios';

const TOKEN_KEY = 'admin_token';

export const getToken = () => localStorage.getItem(TOKEN_KEY);
export const setToken = (token: string) => localStorage.setItem(TOKEN_KEY, token);
export const clearToken = () => localStorage.removeItem(TOKEN_KEY);

const http = axios.create({
  baseURL: process.env.API_BASE,
  timeout: 15000
});

http.interceptors.request.use(config => {
  const token = getToken();
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`);
  }
  return config;
});

http.interceptors.response.use(
  res => res,
  (err: AxiosError<{ error?: ApiError }>) => {
    // 改密码时原密码错误也是 401（INVALID_CREDENTIALS），不能当作登录失效
    if (err.response?.data?.error?.code === 'UNAUTHORIZED' && getToken()) {
      clearToken();
      window.location.replace(`${process.env.PUBLIC_PATH}login`);
    }
    return Promise.reject(err);
  }
);

export interface ApiError {
  status: number;
  code: string;
  message: string;
}

export const toApiError = (err: unknown): ApiError => {
  const e = err as AxiosError<{ error?: ApiError }>;
  const body = e.response?.data?.error;
  return {
    status: e.response?.status ?? 0,
    code: body?.code ?? 'NETWORK',
    message: body?.message ?? e.message ?? 'network error'
  };
};

export interface Page<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface NamedRef {
  id: number;
  name: string;
}

export type ArticleStatus = 'draft' | 'published';

export interface AdminArticle {
  id: number;
  title: string;
  summary: string;
  content?: string;
  status: ArticleStatus;
  category: NamedRef | null;
  tags: NamedRef[];
  publishedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface ArticleInput {
  title: string;
  content: string;
  categoryId: number | null;
  tagIds: number[];
  status: ArticleStatus;
  publishedAt: string;
}

export interface ArticleQuery {
  status: ArticleStatus;
  keyword?: string;
  categoryId?: number;
  tagId?: number;
  page: number;
  pageSize: number;
}

export interface Category extends NamedRef {
  articleCount: number;
}

export interface CategoryList {
  items: Category[];
  uncategorizedCount: number;
}

export type Tag = Category;

export interface AdminComment {
  id: number;
  parentId: number | null;
  article: { id: number; title: string } | null;
  nickname: string;
  email: string;
  website: string;
  avatar: string;
  content: string;
  isAdmin: boolean;
  ip: string;
  createdAt: string;
}

export interface Moment {
  id: number;
  content: string;
  images: string[];
  createdAt: string;
}

export interface MomentInput {
  content: string;
  images: string[];
  createdAt: string;
}

export interface FriendLink {
  id: number;
  name: string;
  url: string;
  avatar: string;
  description: string;
}

export type FriendLinkInput = Omit<FriendLink, 'id'>;

export interface Changelog {
  id: number;
  items: string[];
  loggedAt: string;
}

export type ChangelogInput = Omit<Changelog, 'id'>;

export interface Project {
  id: number;
  name: string;
  description: string;
  cover: string;
  url: string;
  sortOrder: number;
}

export type ProjectInput = Omit<Project, 'id'>;

export type PageKey = 'about-site' | 'about-me';

export interface Stats {
  publishedCount: number;
  draftCount: number;
  categoryCount: number;
  tagCount: number;
  commentCount: number;
  momentCount: number;
  friendLinkCount: number;
  projectCount: number;
  viewCount: number;
}

export interface User {
  id: number;
  email: string;
  nickname: string;
  avatar: string;
  website: string;
}

interface UploadTicket {
  method: string;
  uploadUrl: string;
  headers: Record<string, string>;
  url: string;
}

const get = <T>(url: string, params?: object) =>
  http.get<T>(url, { params }).then(res => res.data);
const post = <T>(url: string, body?: unknown) =>
  http.post<T>(url, body).then(res => res.data);
const put = (url: string, body: unknown) => http.put(url, body).then(() => undefined);
const del = (url: string) => http.delete(url).then(() => undefined);

const resource = <I>(path: string) => ({
  create: (input: I) => post<{ id: number }>(path, input),
  update: (id: number, input: I) => put(`${path}/${id}`, input),
  remove: (id: number) => del(`${path}/${id}`)
});

export const login = (email: string, password: string) =>
  post<{ token: string; user: User }>('/auth/login', { email, password }).then(res => {
    setToken(res.token);
    return res.user;
  });

export const getMe = () => get<User>('/auth/me');

export type ProfileInput = Pick<User, 'nickname' | 'avatar' | 'website'>;

export const updateProfile = (input: ProfileInput) =>
  http.put<User>('/auth/profile', input).then(res => res.data);

// 服务端会让其它设备上的旧 token 失效，并为当前会话签发新 token
export const changePassword = (currentPassword: string, newPassword: string) =>
  http
    .put<{ token: string }>('/auth/password', { currentPassword, newPassword })
    .then(res => setToken(res.data.token));

export const getStats = () => get<Stats>('/admin/stats');

export const articleApi = {
  ...resource<ArticleInput>('/admin/articles'),
  list: (query: ArticleQuery) => get<Page<AdminArticle>>('/admin/articles', query),
  get: (id: number) => get<AdminArticle>(`/admin/articles/${id}`)
};

export const categoryApi = {
  ...resource<{ name: string }>('/admin/categories'),
  list: () => get<CategoryList>('/admin/categories')
};

export const tagApi = {
  ...resource<{ name: string }>('/admin/tags'),
  list: () => get<{ items: Tag[] }>('/admin/tags').then(res => res.items)
};

export const commentApi = {
  list: (page: number, pageSize: number) =>
    get<Page<AdminComment>>('/admin/comments', { page, pageSize }),
  remove: (id: number) => del(`/admin/comments/${id}`)
};

export const momentApi = {
  ...resource<MomentInput>('/admin/moments'),
  list: (page: number, pageSize: number) =>
    get<Page<Moment>>('/admin/moments', { page, pageSize })
};

export const friendLinkApi = {
  ...resource<FriendLinkInput>('/admin/friend-links'),
  list: () => get<{ items: FriendLink[] }>('/admin/friend-links').then(res => res.items)
};

export const changelogApi = {
  ...resource<ChangelogInput>('/admin/changelogs'),
  list: () => get<{ items: Changelog[] }>('/admin/changelogs').then(res => res.items)
};

export const projectApi = {
  ...resource<ProjectInput>('/admin/projects'),
  list: () => get<{ items: Project[] }>('/admin/projects').then(res => res.items)
};

export const pageApi = {
  get: (key: PageKey) => get<{ content: string }>(`/pages/${key}`).then(res => res.content),
  update: (key: PageKey, content: string) => put(`/admin/pages/${key}`, { content })
};

export const noticeApi = {
  get: () => get<{ notice: string }>('/site').then(res => res.notice),
  update: (notice: string) => put('/admin/site/notice', { notice })
};

// 先向服务端申请预签名地址，再由浏览器直传对象存储，返回图片公网地址
export const uploadImage = async (file: File) => {
  const ticket = await post<UploadTicket>('/admin/uploads', {
    filename: file.name,
    size: file.size
  });
  const res = await fetch(ticket.uploadUrl, {
    method: ticket.method,
    headers: ticket.headers,
    body: file
  });
  if (!res.ok) {
    throw new Error(`upload failed: ${res.status}`);
  }
  return ticket.url;
};
