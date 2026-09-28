// 站点个性化配置：换成自己的信息即可，留空的项目不会在页面上显示。
// 图片可以填完整 URL（比如后台上传到 OSS 后得到的链接），也可以放进 public/assets 用 /assets/xxx 引用。
export const siteConfig = {
  // 浏览器标题、首页大标题、文章版权信息
  title: '飞鸟小站',
  // 建站日期，首页「运行天数」从这天开始算
  startTime: '2026-09-28 00:00:00',
  // 本站地址，留言页「本站信息」里展示；留空则用当前访问的域名
  url: '',
  // 后台管理地址，导航栏齿轮按钮跳转；默认同域名 /admin/，构建时可用环境变量 BLOG_ADMIN_URL 覆盖
  adminUrl: process.env.BLOG_ADMIN_URL as string,
  // 页脚「源代码」链接
  sourceUrl: '',
  // ICP 备案号，例如「浙ICP备xxxxxxxx号-1」
  icp: {
    no: '',
    url: 'https://beian.miit.gov.cn/'
  },

  author: {
    name: '飞鸟',
    descr: '一只平凡的鸟罢了。',
    avatar: '/assets/avatar.svg'
  },

  // 首页侧栏社交按钮：链接类填主页地址，二维码类填图片地址
  social: {
    github: '',
    csdn: '',
    weChatQRCode: '',
    qqQRCode: ''
  },

  // 首页侧栏博客卡片上的头像
  cardImage: '/assets/avatar.svg',
  // 友链、评论头像加载前的占位图
  loadingImage: '/assets/loading.svg'
};
