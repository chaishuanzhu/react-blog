import classNames from 'classnames';
import React from 'react';

import { siteConfig } from '@/site.config';
import { useLazyImg } from '@/utils/hooks/useLazyImg';

import s from './index.scss';

interface Props {
  link?: string;
  avatar?: string;
  name?: string;
  descr?: string;
}

const LinkItem: React.FC<Props> = ({ link, avatar, name, descr }) => {
  const { imgRef, imgUrl } = useLazyImg(avatar!, siteConfig.loadingImage);

  return (
    <div className={s.item}>
      <a href={link} rel='noreferrer' target='_blank' className={s.link}>
        <div ref={imgRef} className={s.left}>
          <img
            src={imgUrl}
            className={classNames({
              [s.avatar]: imgUrl !== siteConfig.loadingImage,
              [s.loading]: imgUrl === siteConfig.loadingImage
            })}
          />
        </div>
        <div className={s.right}>
          <div className={s.title}>{name}</div>
          <div className={s.descr}>{descr}</div>
        </div>
      </a>
    </div>
  );
};

export default LinkItem;
