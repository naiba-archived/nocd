# NoCD 持续交付系统

[![GolangCI](https://golangci.com/badges/github.com/naiba/nocd.svg)](https://golangci.com/r/github.com/naiba/nocd) ![构建状态](https://github.com/naiba/nocd/workflows/Build%20Docker%20Image/badge.svg)

**NoCD** 是一个 Go 实现的轻便可控的持续交付系统。

## 界面预览

| ![首页截图](https://github.com/naiba/nocd/raw/master/README/首页截图.png) | ![服务器管理](https://github.com/naiba/nocd/raw/master/README/服务器管理.png) | ![项目管理](https://github.com/naiba/nocd/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://github.com/naiba/nocd/raw/master/README/交付记录.png) | ![管理中心](https://github.com/naiba/nocd/raw/master/README/查看日志.png) | ![查看日志](https://github.com/naiba/nocd/raw/master/README/管理中心.png)  |

## 功能特色

- 服务器：可以添加多个部署服务器。
- 项目：支持解析 Gogs、GitHub、Gitlab、BitBucket 的 WebHook
- 通知：部署成功或失败经 `Server酱` 推送到您的微信
- 交付记录：可以查看部署记录，用户可以停止部署中的流程
- 管理面板：查看系统状态，管理用户，管理部署中的流程

## 部署指北

### Docker

1. 创建配置文件（如`/data/nocd`文件夹）

   ```shell
   nano /data/nocd/app.ini
   ```

   文件内容参考下面（ `web_listen = 0.0.0.0:8000` 配置不要改）

2. 运行NoCD

   ```
   docker run -d --name=nocd -p 8000:8000 -v /data/nocd/:/data/conf docker.pkg.github.com/naiba/nocd/app:latest
   ```

### 源码编译

1. Clone 源代码

2. 进入应用目录 `cd nocd/cmd/web`

3. 打包资源文件并编译

   ```shell
   go get -u github.com/tmthrgd/go-bindata/go-bindata
   go-bindata resource/...
   go build
   ```

4. 在 `conf/app.ini` 创建配置文件

   ```ini
   [nocd]
   cookie_key_pair = example
   debug = true
   domain = cd.git.cm
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

6. 在 `GitHub` 设置回调：`https://cd.git.cm/oauth2/callback`

## 常见问题

1. 为什么我的部署脚本总是执行失败 或者 根本没有执行？

    > 请检查您的 PATH 路径是否引入，建议提前 export 一下路径，自动部署的时候不会
    >
    > `source .bash_profile`。

2. 如何保持后台运行？<br>

    > 可以使用` systemd` 。 更推荐使用docker方式运行。


## 版权声明

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnaiba%2Fnocd.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnaiba%2Fnocd?ref=badge_large)

Copy &copy; 2017-2019 Naiba
