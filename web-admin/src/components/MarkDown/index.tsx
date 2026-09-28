import './hljs.custom.scss';

import classNames from 'classnames';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';
import { Marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
import React, { useMemo } from 'react';

import s from './index.scss';

interface Props {
  content: string;
  className?: string;
  style?: React.CSSProperties;
  onScroll?: React.UIEventHandler<HTMLDivElement>;
  ref?: React.Ref<HTMLDivElement>;
}

const marked = new Marked(
  markedHighlight({
    highlight: (code, lang) =>
      lang && hljs.getLanguage(lang)
        ? hljs.highlight(code, { language: lang }).value
        : hljs.highlightAuto(code).value
  }),
  { gfm: true, breaks: true }
);

const MarkDown: React.FC<Props> = ({ content = '', className, onScroll, style, ref }) => {
  const html = useMemo(() => {
    const raw = marked.parse(content, { async: false }).replace(/<pre>/g, "<pre id='hljs'>");
    return DOMPurify.sanitize(raw);
  }, [content]);

  return (
    <div
      style={style}
      ref={ref}
      onScroll={onScroll}
      className={classNames(s.markdownBox, className)}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
};

export default MarkDown;
