## NoCD Continuous Delivery System

![Build Status](https://github.com/naiba/nocd/workflows/Build%20Docker%20Image/badge.svg)

**NoCD** is a lightweight and controllable continuous delivery system implemented by Go.

## Interface preview

| ![首页截图](https://github.com/naiba/nocd/raw/master/README/首页截图.png) | ![服务器管理](https://github.com/naiba/nocd/raw/master/README/服务器管理.png) | ![项目管理](https://github.com/naiba/nocd/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://github.com/naiba/nocd/raw/master/README/交付记录.png) | ![管理中心](https://github.com/naiba/nocd/raw/master/README/查看日志.png) | ![查看日志](https://github.com/naiba/nocd/raw/master/README/管理中心.png)  |

## Features

- Multi Language support: English, Chinese (PR is welcome)
- Server: Multiple deployment servers can be added
- Project: Support parsing Webhooks of various popular Git hosting platforms
- Notification: Flexible custom Webhook
- Delivery record: You can view the deployment record, and the user can stop the deployment process
- Management panel: View system status, manage users, and manage deployment processes

## Deployment means north

### Docker

1. Create a configuration file (eg `/data/nocd` folder)

   ```shell
   nano /data/nocd/app.ini
   ```

   Refer to the following for the content of the file (`web_listen = 0.0.0.0:8000` configuration do not change)

2. Run NoCD

   ```shell
   docker run -d --name=nocd -p 8000:8000 -v /data/nocd/:/data/conf docker.pkg.github.com/naiba/dockerfiles/nocd:latest
   ```

### Source code compilation

1. Clone source code

2. Enter the application directory `cd nocd/cmd/web`

3. Compile the binary

   ```shell
   go build
   ```

4. Create a configuration file in `conf/app.ini`

   ```ini
   [nocd]
   cookie_key_pair = i_love_NoCD
   debug = true
   domain = your_domain_name # or ip:port
   web_listen = 0.0.0.0:8000
   loc = Asia/Shanghai
   [third_party]
   github_oauth2_client_id = example
   github_oauth2_client_secret = example
   google_analysis = "NB-XXXXXX-1" # optional
   sentry_dsn = "https://example:xx@example.io/project_id" # optional
   ```

5. Run

   ```shell
   ./web
   ```

6. Set the callback in `GitHub`: `http(s)://your_domain_name/oauth2/callback`

## Common issue

1. Why does my deployment script always fail to execute or not executed at all?

    > Please check whether your PATH path is imported, it is recommended to export the path in advance, it will not be automatically deployed
    >
    > `source .bash_profile`.

2. How to keep running in the background?

    > You can use `systemd`. It is more recommended to run in docker mode.

## License

MIT