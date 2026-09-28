import { useMount, useSafeState } from 'ahooks';
import dayjs from 'dayjs';

import { siteConfig } from '@/site.config';

export const useRunTime = () => {
  const [runTime, setRunTime] = useSafeState(0);

  useMount(() => {
    const nowTime = new Date().getTime();
    const startTime = new Date(siteConfig.startTime).getTime();
    const runTime = dayjs(nowTime).diff(dayjs(startTime), 'days');
    setRunTime(runTime);
  });

  return { runTime };
};
