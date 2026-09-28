import './hljs.custom.scss';

import classNames from 'classnames';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';
import { Marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
import React from 'react';

import s from './index.scss';

interface Props {
  content?: string;
  className?: string;
}

hljs.configure({
  classPrefix: 'hljs-',
  languages: ['CSS', 'HTML', 'JavaScript', 'TypeScript', 'Markdown']
});

const marked = new Marked(
  markedHighlight({
    highlight: (code, lang) =>
      lang && hljs.getLanguage(lang)
        ? hljs.highlight(code, { language: lang }).value
        : hljs.highlightAuto(code).value
  }),
  {
    gfm: true, // 允许 GitHub 标准的 markdown
    breaks: true // 允许回车换行，要求 gfm 为 true
  }
);

const MarkDown: React.FC<Props> = ({ content, className }) => {
  const html = marked.parse(content || '', { async: false }).replace(/<pre>/g, "<pre id='hljs'>");

  return (
    <div
      className={classNames(s.marked, className)}
      dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(html) }}
    />
  );
};

export default MarkDown;
