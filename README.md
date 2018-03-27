# GoCD 持续交付系统

[![Go Report Card](https://goreportcard.com/badge/git.cm/naiba/gocd)](https://goreportcard.com/report/git.cm/naiba/gocd)  [![Build status](https://ci.appveyor.com/api/projects/status/d7bo0ng4n0bm8l11?svg=true)](https://ci.appveyor.com/project/naiba/gocd)  [![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)


**GoCD** 是一个 Golang 实现的持续交付系统。

## 界面预览

| ![首页截图](https://git.cm/naiba/gocd/raw/master/README/首页截图.png) | ![服务器管理](https://git.cm/naiba/gocd/raw/master/README/服务器管理.png) | ![项目管理](https://git.cm/naiba/gocd/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://git.cm/naiba/gocd/raw/master/README/交付记录.png) | ![管理中心](https://git.cm/naiba/gocd/raw/master/README/查看日志.png) | ![查看日志](https://git.cm/naiba/gocd/raw/master/README/管理中心.png)  |

## 部署教程

1. Clone 源代码

2. 进入应用目录 `gocd/cmd/web`

3. 打包资源文件并编译

       ```shell
   go get -u github.com/tmthrgd/go-bindata/...
   go-bindata resource/...
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
   google_analysis = "NB-XXXXXX-1"
   [third_party]
   github_oauth2_client_id = example
   github_oauth2_client_secret = example
   sentry_dsn = "https://example:xx@example.io/"
       ```

5. 运行

       ```shell
   ./web
       ```

6. 在 `GitHub` 设置回调：`http://mjj.cx/oauth2/callback`

## FAQ

1. 为什么我的部署脚本总是执行失败 或者 根本没有执行？<br>
  `请检查您的 PATH 路径是否引入，建议提前 export 一下路径，自动部署的时候不会 source .bash_profile 。`

2. 如何保持后台运行？<br>
  `可以使用 systemd 。`


## 版权声明

本仓库代码遵循 MIT 协议

Copy &copy; 2018 Naiba