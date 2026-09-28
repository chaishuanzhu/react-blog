/**
 * 生成指定范围内的随机整数，左闭右闭
 * @param {Number} Min
 * @param {Number} Max
 * @return {Number}
 */
export const getRandomNum = (Min: number, Max: number) => {
  const Range = Max - Min + 1;
  const Rand = Math.random();
  return Min + Math.floor(Rand * Range);
};

/**
 * 打乱数组
 * @param {T[]} array
 * @return {T[]}
 */
export const shuffleArray = <T>(array: T[]): T[] => {
  if (!array) return [];
  const res = [...array];
  const len = res.length;
  for (let i = len - 1; i > 0; i--) {
    const randomPos = Math.floor(Math.random() * (i + 1));
    [res[i], res[randomPos]] = [res[randomPos], res[i]];
  }
  return res;
};
