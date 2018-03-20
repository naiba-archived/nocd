# GoCD 持续交付系统

**GoCD** 是一个 Golang 实现的持续交付系统。

## 部署教程

1. Clone 源代码到本地

2. 在 `cmd/web` 中执行安装所需第三方包

   ```shell
   go get
   ```

3. 在 `cmd/web` 中执行编译为二进制文件

   ```shell
   go build
   ```

4. 在 `cmd/web/conf/app.ini` 创建配置文件

   ```ini
   [gocd]
   debug = true
   loc = Asia/Shanghai
   cookie_key_pair = [自定义密钥]
   domain = 0.0.0.0:8000
   [third_party]
   sentry_dsn = [Sentry服务DSN]
   github_oauth2_client_id = [GitHub client ID]
   github_oauth2_client_secret = [GitHub client secret]
   ```

5. 运行

   ```shell
   ./web
   ```

## 界面预览

| ![首页截图](https://git.cm/naiba/GoCD/raw/master/README/首页截图.png) | ![服务器管理](https://git.cm/naiba/GoCD/raw/master/README/服务器管理.png) | ![项目管理](https://git.cm/naiba/GoCD/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://git.cm/naiba/GoCD/raw/master/README/交付记录.png) | ![查看日志](https://git.cm/naiba/GoCD/raw/master/README/查看日志.png) |                                                              |

## 版权声明

本系统使用 MIT 协议

Copy &copy; 2018 Naiba