import dayjs from 'dayjs';
import React from 'react';

import { siteConfig } from '@/site.config';

import s from './index.scss';

interface Props {
  content: string;
  date: string;
  imgs: string[];
  handlePreView: (url: string) => void;
}

const SayPop: React.FC<Props> = ({ content, date, imgs, handlePreView }) => (
  <div className={s.sayItem}>
    <div className={s.avatarBox}>
      <img src={siteConfig.author.avatar} className={s.avatar} />
    </div>

    <div className={s.contentBox}>
      <div className={s.content}>
        {content}
        <span className={s.date}>{dayjs(date).format('YYYY-MM-DD HH:mm:ss')}</span>
        {!!imgs?.length && (
          <div className={s.sayImgsBox}>
            {imgs.map((img, index) => (
              <div key={index} className={s.sayImg} onClick={() => handlePreView(img)}>
                <img src={img} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  </div>
);

export default SayPop;
