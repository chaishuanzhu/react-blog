import { useThrottleFn } from 'ahooks';
import classNames from 'classnames';
import React, { useEffect, useState } from 'react';

interface Heading {
  el: HTMLElement;
  text: string;
  level: number;
}

interface Props {
  // 渲染后的文章容器选择器
  target: string;
  // 变化时重新收集标题
  content?: string;
  headingTopOffset: number;
  className?: string;
  onNavItemClick?: () => void;
}

const collectHeadings = (target: string): Heading[] => {
  const root = document.querySelector(target);
  const els = Array.from(root?.querySelectorAll<HTMLElement>('h1, h2, h3, h4, h5, h6') ?? []);
  const depths = els.map(el => Number(el.tagName[1]));
  const minDepth = Math.min(...depths);
  // 最高一级标题对应 title-level2，与原有样式保持一致
  return els.map((el, i) => ({ el, text: el.textContent?.trim() || '', level: depths[i] - minDepth + 2 }));
};

const Toc: React.FC<Props> = ({ target, content, headingTopOffset, className, onNavItemClick }) => {
  const [headings, setHeadings] = useState<Heading[]>([]);
  const [active, setActive] = useState(-1);

  useEffect(() => {
    setHeadings(collectHeadings(target));
  }, [target, content]);

  const { run: onScroll } = useThrottleFn(
    () => {
      let current = -1;
      headings.forEach((h, i) => {
        if (h.el.getBoundingClientRect().top <= headingTopOffset + 1) current = i;
      });
      setActive(current);
    },
    { wait: 100 }
  );

  useEffect(() => {
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, [headings, onScroll]);

  const scrollTo = (h: Heading) => {
    const top = h.el.getBoundingClientRect().top + window.scrollY - headingTopOffset;
    window.scrollTo({ top, behavior: 'smooth' });
    onNavItemClick?.();
  };

  if (!headings.length) return null;

  return (
    <div className={classNames('markdown-navigation', className)}>
      {headings.map((h, i) => (
        <div
          key={i}
          className={classNames('title-anchor', `title-level${h.level}`, { active: i === active })}
          onClick={() => scrollTo(h)}
        >
          {h.text}
        </div>
      ))}
    </div>
  );
};

export default Toc;
