// 一些公共常量
const path = require('path');

const ROOT_PATH = path.resolve(__dirname, '../');

const SERVER_HOST = 'localhost';
const SERVER_PORT = 3001;

// 开发时 /api 请求代理到 Go 服务端
const API_TARGET = process.env.API_TARGET || 'http://localhost:8080';

// 后台挂载的路径前缀，须与服务端 internal/http/router.go 的 adminPrefix 一致
const PUBLIC_PATH = '/admin/';

module.exports = {
  ROOT_PATH,
  SERVER_HOST,
  SERVER_PORT,
  API_TARGET,
  PUBLIC_PATH
};
