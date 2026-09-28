import { siteConfig } from '@/site.config';

const absolute = (url: string) => new URL(url, window.location.origin).href;

export const useSite = () => {
  const { url, author } = siteConfig;
  const mySite = [
    {
      key: 'name',
      value: author.name
    },
    {
      key: 'link',
      value: url || window.location.origin
    },
    {
      key: 'avatar',
      value: absolute(author.avatar)
    },
    {
      key: 'descr',
      value: author.descr
    }
  ];

  return { mySite };
};
