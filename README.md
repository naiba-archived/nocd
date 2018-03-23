# GoCD 持续交付系统

[![Build status](https://ci.appveyor.com/api/projects/status/d7bo0ng4n0bm8l11?svg=true)](https://ci.appveyor.com/project/naiba/gocd)

**GoCD** 是一个 Golang 实现的持续交付系统。

## 部署教程

1. Clone 源代码到本地

2. 进入应用目录 `gocd/cmd/web`

2. 安装依赖

       ```shell
       go get
       ```

3. 编译

       ```shell
       go build
       ```

4. 在 `conf/app.ini` 创建配置文件

       ```ini
       [gocd]
       cookie_key_pair = example
       debug = true
       domain = mjj.cx
       web_listen = 0.0.0.0:8000
       loc = Asia/Shanghai
       [third_party]
       github_oauth2_client_id = example
       github_oauth2_client_secret = example
       sentry_dsn = "https://example:xx@example.io/"
       ```

5. 运行

       ```shell
       ./web
       ```

6. `GitHub` 回调地址：`http://mjj.cx/oauth2/callback`

## 界面预览

| ![首页截图](https://git.cm/naiba/gocd/raw/master/README/首页截图.png) | ![服务器管理](https://git.cm/naiba/gocd/raw/master/README/服务器管理.png) | ![项目管理](https://git.cm/naiba/gocd/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://git.cm/naiba/gocd/raw/master/README/交付记录.png) | ![查看日志](https://git.cm/naiba/gocd/raw/master/README/查看日志.png) |                                                              |

## 版权声明

本仓库代码遵循 MIT 协议

Copy &copy; 2018 Naiba