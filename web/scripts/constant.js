// 一些公共常量
const path = require('path');

const ROOT_PATH = path.resolve(__dirname, '../');

const SERVER_HOST = 'localhost';
const SERVER_PORT = 3000;

// 开发时 /api 请求代理到 Go 服务端
const API_TARGET = process.env.API_TARGET || 'http://localhost:8080';

module.exports = {
  ROOT_PATH,
  SERVER_HOST,
  SERVER_PORT,
  API_TARGET
};
