# React Blog

一个前后端完整的个人博客系统：博客前台、后台管理和 Go 后端放在同一个仓库，打包成一个 Docker 容器，用同一个域名访问。

本项目基于 [lzxjack/react-blog](https://github.com/lzxjack/react-blog) 和 [lzxjack/react-blog-admin](https://github.com/lzxjack/react-blog-admin) 二次开发：把原来依赖的腾讯云 CloudBase 换成了自建的 Go + MySQL 后端，并重写了两个前端的数据层。感谢原作者飞鸟。

## 功能

- **博客前台**：文章列表与详情（Markdown、代码高亮、目录导航）、分类 / 标签 / 关键字搜索、说说、留言板与文章评论、友链、建站日志、关于页、三种主题切换
- **后台管理**：文章与草稿、分类与标签、说说、评论、友链、建站日志、关于页、站点公告的增删改查，图片上传到 S3 兼容存储，数据概览
- **后端**：REST API、JWT 登录（改密码后旧 token 自动失效）、评论限流、评论与回复的邮件通知、启动时自动执行数据库迁移、Swagger UI 接口文档、CloudBase 数据导入工具

## 技术栈

| 目录 | 说明 | 主要技术 |
| --- | --- | --- |
| `server/` | 后端 API，同时托管两个前端的静态文件 | Go 1.26、Gin、MySQL 8、sqlc、golang-migrate |
| `web/` | 博客前台，访问路径 `/` | React 19、React Router、Ant Design 6、webpack 5 |
| `web-admin/` | 后台管理，访问路径 `/admin/` | React 19、React Router、Arco Design、webpack 5 |

## 本地开发

需要 Go 1.26+、Node.js 24+、MySQL 8。

1. 准备数据库和配置：

   ```bash
   cp server/.env.example server/.env
   # 修改 DATABASE_DSN、JWT_SECRET、ADMIN_EMAIL、ADMIN_PASSWORD 等
   ```

   本机没有 MySQL 的话，可以用 compose 里的可选服务启动一个：`docker compose --profile mysql up -d mysql`。

2. 启动后端（默认监听 `:8080`，首次启动会建表并按 `ADMIN_EMAIL` / `ADMIN_PASSWORD` 创建管理员）：

   ```bash
   cd server && make run
   ```

3. 分别启动两个前端，开发服务器会把 `/api` 代理到 `http://localhost:8080`（可用环境变量 `API_TARGET` 修改）：

   ```bash
   cd web && npm install && npm start             # http://localhost:3000
   cd web-admin && npm install && npm start       # http://localhost:3001/admin/
   ```

开发环境下接口文档在 <http://localhost:8080/api/docs/>，规范文件是 [`server/api/openapi.yaml`](server/api/openapi.yaml)。

常用命令（在 `server/` 下执行）：`make test` 跑测试，`make lint` 做静态检查，`make sqlc` 修改 SQL 后重新生成代码。

## 部署

整个项目构建成一个镜像：Go 服务在 8080 端口同时提供博客前台（`/`）、后台（`/admin/`）和 API（`/api/v1`）。

```bash
cp server/.env.example server/.env   # 填好生产配置
docker compose up -d --build
```

生产环境需要注意：

- 设置 `SITE_URL` 为博客的公网地址（邮件通知里的链接会用到），并换掉 `JWT_SECRET`、`ADMIN_PASSWORD`
- 容器只提供 HTTP，请在前面加一层 HTTPS（CDN、云负载均衡，或宿主机上的 nginx / Caddy）；代理时要保留 `Host` 请求头，并把代理地址填进 `TRUSTED_PROXIES`，评论限流才能拿到访客真实 IP
- `APP_ENV=production` 时默认关闭接口文档，需要的话设置 `API_DOCS=true`
- 图片上传使用 S3 兼容存储（MinIO、RustFS、AWS S3、R2 等），建议使用只能写入该 bucket 的密钥，bucket 需要允许公开读取，并配置允许后台域名 PUT 的 CORS
- 不配置 `SMTP_HOST` 则不发送邮件通知

全部配置项及说明见 [`server/.env.example`](server/.env.example)。

## 个性化

站点标题、作者信息、社交链接、备案号等在这两个文件里修改：

- 前台：[`web/src/site.config.ts`](web/src/site.config.ts)
- 后台：[`web-admin/src/site.config.ts`](web-admin/src/site.config.ts)

## 从 CloudBase 迁移数据

如果之前用的是原版（CloudBase）博客，可以在云开发控制台把各个集合导出为 JSON，放进同一个目录（文件名与集合名一致，如 `articles.json`、`allComments.json`），然后导入：

```bash
cd server
make import DIR=./cloudbase-export ARGS="-dry-run"                                  # 先试运行，只报告不写入
make import DIR=./cloudbase-export ARGS="-truncate -admin-email you@example.com"   # 清空现有内容后导入
```

`-admin-email` 指定的邮箱发表的评论会标记为博主回复。

## 许可证

[MIT](LICENSE)
