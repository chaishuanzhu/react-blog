const path = require('path');
const webpack = require('webpack');
const WebpackBar = require('webpackbar');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

const { ROOT_PATH, PUBLIC_PATH } = require('../constant');
const { isDevelopment, isProduction } = require('../env');

// modules 为 true 时启用 CSS Modules（默认导出类名映射）
const getCssLoaders = (modules = true) => {
  const cssLoaders = [
    isDevelopment
      ? 'style-loader'
      : { loader: MiniCssExtractPlugin.loader, options: { publicPath: '../' } },
    {
      loader: 'css-loader',
      options: {
        modules: modules
          ? {
              // 模块化类名，防止重复
              localIdentName: '[local]--[hash:base64:10]',
              namedExport: false,
              exportLocalsConvention: 'as-is'
            }
          : false,
        sourceMap: isDevelopment
      }
    }
  ];

  // 生产模式时，才需要加css前缀
  isProduction &&
    cssLoaders.push({
      loader: 'postcss-loader',
      options: {
        postcssOptions: {
          plugins: [['postcss-preset-env', { autoprefixer: { grid: true } }]]
        }
      }
    });

  return cssLoaders;
};

const sassLoader = {
  loader: 'sass-loader',
  options: { sourceMap: isDevelopment }
};

module.exports = {
  entry: {
    index: path.resolve(ROOT_PATH, './src/index')
  },

  output: {
    // 自动删除上一次打包的产物
    clean: true
  },

  plugins: [
    // html模板
    new HtmlWebpackPlugin({
      template: path.resolve(ROOT_PATH, './public/index.html'),
      filename: 'index.html',
      inject: 'body',
      templateParameters: { PUBLIC_PATH }
    }),
    // 打包显示进度条
    new WebpackBar(),
    // webpack打包不会有类型检查，强制ts类型检查
    new ForkTsCheckerWebpackPlugin({
      typescript: {
        configFile: path.resolve(ROOT_PATH, './tsconfig.json')
      }
    }),
    // 复制不用动态导入的资源
    new CopyWebpackPlugin({
      patterns: [
        {
          context: 'public',
          from: 'assets/*',
          to: path.resolve(ROOT_PATH, './build'),
          toType: 'dir',
          noErrorOnMissing: true,
          globOptions: {
            dot: true,
            ignore: ['**/index.html'] // **表示任意目录下
          }
        }
      ]
    }),
    // 接口地址，默认与前端同源（开发时由 devServer 代理，生产由网关转发）
    new webpack.DefinePlugin({
      'process.env.API_BASE': JSON.stringify(process.env.API_BASE || '/api/v1'),
      'process.env.PUBLIC_PATH': JSON.stringify(PUBLIC_PATH),
      // 博客前台地址；生产环境与后台同域名部署，默认就是根路径
      'process.env.BLOG_URL': JSON.stringify(
        process.env.BLOG_URL || (isDevelopment ? 'http://localhost:3000' : '/')
      )
    })
  ],

  module: {
    rules: [
      {
        test: /\.css$/,
        include: /node_modules/,
        use: getCssLoaders(false)
      },
      {
        test: /\.css$/,
        exclude: /node_modules/,
        use: getCssLoaders()
      },
      {
        test: /\.scss$/,
        exclude: [/node_modules/, /\.custom.scss$/],
        use: [...getCssLoaders(), sassLoader]
      },
      {
        test: /\.custom.scss$/,
        use: [...getCssLoaders(false), sassLoader]
      },
      {
        test: /\.(tsx?|js)$/, // ts\tsx\js
        loader: 'babel-loader',
        options: { cacheDirectory: true }, // 缓存公共文件
        exclude: /node_modules/
      },
      {
        test: [/\.bmp$/, /\.gif$/, /\.jpe?g$/, /\.png$/],
        // 自动选择导出为单独文件还是url形式
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 4 * 1024
          }
        }
      },
      {
        test: /\.(eot|svg|ttf|woff|woff2?)$/,
        // 分割为单独文件，并导出url
        type: 'asset/resource'
      }
    ]
  },

  // 路径配置别名
  resolve: {
    alias: {
      '@': path.resolve(ROOT_PATH, './src')
    },
    // 若没有写后缀时，依次从数组中查找相应后缀文件是否存在
    extensions: ['.tsx', '.ts', '.js', '.json'],
    fallback: { crypto: false }
  },

  // 缓存
  cache: {
    // 基于文件系统的持久化缓存
    type: 'filesystem',
    buildDependencies: {
      // 当配置文件发生变化时，缓存失效
      config: [__filename]
    }
  }
};
