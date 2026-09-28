import { ArrowRightOutlined, RedoOutlined } from '@ant-design/icons';
import { useKeyPress, useSafeState } from 'ahooks';
import { message } from 'antd';
import React, { useRef } from 'react';

import s from './index.scss';

interface Props {
  page: number;
  setPage: (page: number) => void;
  keyword: string;
  setKeyword: (keyword: string) => void;
}

const Search: React.FC<Props> = ({ page, setPage, keyword, setKeyword }) => {
  const [input, setInput] = useSafeState('');
  const inputRef = useRef(null);

  const search = () => {
    const value = input.trim();
    if (!value) {
      message.info('请输入关键词再搜索!');
      return;
    }
    setKeyword(value);
    setPage(1);
  };

  const reset = () => {
    if (!keyword && page === 1 && !input) {
      message.info('无需重置!');
      return;
    }
    setInput('');
    setKeyword('');
    setPage(1);
  };

  useKeyPress(13, search, {
    target: inputRef
  });

  useKeyPress(27, reset, {
    target: inputRef
  });

  return (
    <div className={s.searchBox}>
      <input
        ref={inputRef}
        autoFocus
        type='text'
        placeholder='搜索文章标题...'
        className={s.search}
        value={input}
        onChange={e => setInput(e.target.value)}
      />
      {/* 搜索按钮 */}
      <div className={s.searchBtn} onClick={search}>
        <ArrowRightOutlined />
      </div>
      {/* 重置按钮 */}
      <div className={s.searchBtn} onClick={reset}>
        <RedoOutlined />
      </div>
    </div>
  );
};

export default Search;
