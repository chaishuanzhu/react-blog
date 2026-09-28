import { Message } from '@arco-design/web-react';
import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat';

import { toApiError } from './api';

dayjs.extend(customParseFormat);

export const showError = (err: unknown) => {
  Message.warning(`操作失败：${toApiError(err).message}`);
};

// 执行写操作并提示结果，成功返回 true
export const mutate = async (action: () => Promise<unknown>, successText: string) => {
  try {
    await action();
    Message.success(successText);
    return true;
  } catch (err) {
    showError(err);
    return false;
  }
};

// 把输入框里的本地时间转成 ISO 字符串，非法时返回 null
export const parseLocalTime = (text: string, format: string) => {
  const d = dayjs(text.trim(), format, true);
  return d.isValid() ? d.toISOString() : null;
};

// 删除一项后，若当前页已空则回退一页
export const pageAfterDelete = (total: number, page: number, pageSize: number) =>
  Math.max(1, Math.min(page, Math.ceil((total - 1) / pageSize)));
