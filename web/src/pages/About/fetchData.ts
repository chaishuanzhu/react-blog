import type { SitePage} from '@/utils/api';
import { getCategories, getPage, toApiError } from '@/utils/api';

const getPageOrEmpty = (key: 'about-site' | 'about-me'): Promise<SitePage> =>
  getPage(key).catch(err => {
    if (toApiError(err).status === 404) return { key, content: '' };
    throw err;
  });

export const fetchData = async () => {
  const [site, me, categories] = await Promise.all([
    getPageOrEmpty('about-site'),
    getPageOrEmpty('about-me'),
    getCategories()
  ]);

  return { site, me, categories };
};
