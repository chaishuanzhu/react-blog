// 后台个性化配置。public/index.html 里的 <title> 是页面加载前的占位标题，一并修改即可。
export const siteConfig = {
  // 浏览器标题
  title: '飞鸟小站后台管理',
  // 登录页头像，以及账号未设置头像时顶栏显示的默认头像
  defaultAvatar: `${process.env.PUBLIC_PATH}assets/avatar.svg`,
  // 博客前台地址，顶栏主页按钮跳转；默认同域名根路径，构建时可用环境变量 BLOG_URL 覆盖
  blogUrl: process.env.BLOG_URL as string
};
