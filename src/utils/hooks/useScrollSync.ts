import { useThrottleFn } from 'ahooks';
import type { UIEvent } from 'react';
import { useRef } from 'react';

// 编辑区与预览区按滚动比例同步
export const useScrollSync = () => {
  const leftRef = useRef<HTMLTextAreaElement>(null);
  const rightRef = useRef<HTMLDivElement>(null);

  const handleScroll = (event: UIEvent<HTMLElement>) => {
    const source = event.target as HTMLElement;
    const target = source === leftRef.current ? rightRef.current : leftRef.current;
    if (!target) return;
    const ratio = source.scrollTop / (source.scrollHeight - source.clientHeight || 1);
    target.scrollTop = ratio * (target.scrollHeight - target.clientHeight);
  };

  const { run: handleScrollRun } = useThrottleFn(handleScroll, { wait: 60 });

  return { leftRef, rightRef, handleScrollRun };
};
