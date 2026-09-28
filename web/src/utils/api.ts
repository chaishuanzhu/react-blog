import type { AxiosError } from 'axios';
import axios from 'axios';

const TOKEN_KEY = 'token';

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
  (err: AxiosError) => {
    if (err.response?.status === 401) clearToken();
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

export interface ArticleSummary {
  id: number;
  title: string;
  summary: string;
  category: NamedRef | null;
  tags: NamedRef[];
  publishedAt: string;
  updatedAt: string;
}

export interface ArticleDetail extends ArticleSummary {
  content: string;
}

export interface Category extends NamedRef {
  articleCount: number;
}

export interface CategoryList {
  items: Category[];
  uncategorizedCount: number;
}

export interface Tag extends NamedRef {
  articleCount: number;
}

export interface Moment {
  id: number;
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

export interface Changelog {
  id: number;
  items: string[];
  loggedAt: string;
}

export interface Project {
  id: number;
  name: string;
  description: string;
  cover: string;
  url: string;
}

export interface SitePage {
  key: string;
  content: string;
}

export interface Site {
  notice: string;
  viewCount: number;
  articleCount: number;
  categoryCount: number;
  tagCount: number;
}

export interface Comment {
  id: number;
  parentId: number | null;
  nickname: string;
  website: string;
  avatar: string;
  content: string;
  isAdmin: boolean;
  createdAt: string;
}

export interface CommentThread extends Comment {
  replies: Comment[];
}

export interface CommentInput {
  articleId?: number;
  parentId?: number;
  nickname: string;
  email: string;
  website: string;
  content: string;
}

export interface User {
  id: number;
  email: string;
  nickname: string;
  avatar: string;
  website: string;
}

export interface ArticleQuery {
  page?: number;
  pageSize?: number;
  keyword?: string;
  category?: string;
  tag?: string;
}

const get = <T>(url: string, params?: object) =>
  http.get<T>(url, { params }).then(res => res.data);

export const getArticles = (query: ArticleQuery) =>
  get<Page<ArticleSummary>>('/articles', query);

export const getArticle = (id: string) => get<ArticleDetail>(`/articles/${encodeURIComponent(id)}`);

export const getCategories = () => get<CategoryList>('/categories');

export const getTags = () => get<{ items: Tag[] }>('/tags').then(res => res.items);

export const getMoments = (page = 1, pageSize = 100) =>
  get<Page<Moment>>('/moments', { page, pageSize });

export const getFriendLinks = () =>
  get<{ items: FriendLink[] }>('/friend-links').then(res => res.items);

export const getChangelogs = () =>
  get<{ items: Changelog[] }>('/changelogs').then(res => res.items);

export const getProjects = () =>
  get<{ items: Project[] }>('/projects').then(res => res.items);

export const getPage = (key: 'about-site' | 'about-me') => get<SitePage>(`/pages/${key}`);

export const getSite = () => get<Site>('/site');

export const recordView = () =>
  http.post<{ viewCount: number }>('/site/views').then(res => res.data.viewCount);

export const getComments = (articleId: number | undefined, page: number, pageSize: number) =>
  get<Page<CommentThread>>('/comments', { articleId, page, pageSize });

export const postComment = (input: CommentInput) =>
  http.post<Comment>('/comments', input).then(res => res.data);

export const login = (email: string, password: string) =>
  http
    .post<{ token: string; user: User }>('/auth/login', { email, password })
    .then(res => {
      setToken(res.data.token);
      return res.data.user;
    });

export const getMe = () => get<User>('/auth/me');
