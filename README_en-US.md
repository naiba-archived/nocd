## NoCD Continuous Delivery System

![Build Status](https://github.com/naiba/nocd/workflows/Build%20Docker%20Image/badge.svg) <a href="README.md">
    <img height="20px" src="https://img.shields.io/badge/CN-flag.svg?color=555555&style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxMjAwIDgwMCIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiPg0KPHBhdGggZmlsbD0iI2RlMjkxMCIgZD0ibTAsMGgxMjAwdjgwMGgtMTIwMHoiLz4NCjxwYXRoIGZpbGw9IiNmZmRlMDAiIGQ9Im0tMTYuNTc5Niw5OS42MDA3bDIuMzY4Ni04LjEwMzItNi45NTMtNC43ODgzIDguNDM4Ni0uMjUxNCAyLjQwNTMtOC4wOTI0IDIuODQ2Nyw3Ljk0NzkgOC40Mzk2LS4yMTMxLTYuNjc5Miw1LjE2MzQgMi44MTA2LDcuOTYwNy02Ljk3NDctNC43NTY3LTYuNzAyNSw1LjEzMzF6IiB0cmFuc2Zvcm09Im1hdHJpeCg5LjkzMzUyIC4yNzc0NyAtLjI3NzQ3IDkuOTMzNTIgMzI0LjI5MjUgLTY5NS4yNDE1KSIvPg0KPHBhdGggZmlsbD0iI2ZmZGUwMCIgaWQ9InN0YXIiIGQ9Im0zNjUuODU1MiwzMzIuNjg5NWwyOC4zMDY4LDExLjM3NTcgMTkuNjcyMi0yMy4zMTcxLTIuMDcxNiwzMC40MzY3IDI4LjI1NDksMTEuNTA0LTI5LjU4NzIsNy40MzUyLTIuMjA5NywzMC40MjY5LTE2LjIxNDItMjUuODQxNS0yOS42MjA2LDcuMzAwOSAxOS41NjYyLTIzLjQwNjEtMTYuMDk2OC0yNS45MTQ4eiIvPg0KPGcgZmlsbD0iI2ZmZGUwMCI+DQo8cGF0aCBkPSJtNTE5LjA3NzksMTc5LjMxMjlsLTMwLjA1MzQtNS4yNDE4LTE0LjM5NDUsMjYuODk3Ni00LjMwMTctMzAuMjAyMy0zMC4wMjkzLTUuMzc4MSAyNy4zOTQ4LTEzLjQyNDItNC4xNjQ3LTMwLjIyMTUgMjEuMjMyNiwyMS45MDU3IDI3LjQ1NTQtMTMuMjk5OC0xNC4yNzIzLDI2Ljk2MjcgMjEuMTMzMSwyMi4wMDE3eiIvPg0KPHBhdGggZD0ibTQ1NS4yNTkyLDMxNS45Nzk1bDkuMzczNC0yOS4wMzE0LTI0LjYzMjUtMTcuOTk3OCAzMC41MDctLjA1NjYgOS41MDUtMjguOTg4NiA5LjQ4MSwyOC45OTY0IDMwLjUwNywuMDgxOC0yNC42NDc0LDE3Ljk3NzQgOS4zNDkzLDI5LjAzOTItMjQuNzE0LTE3Ljg4NTgtMjQuNzI4OCwxNy44NjUzeiIvPg0KPC9nPg0KPHVzZSB4bGluazpocmVmPSIjc3RhciIgdHJhbnNmb3JtPSJtYXRyaXgoLjk5ODYzIC4wNTIzNCAtLjA1MjM0IC45OTg2MyAxOS40MDAwNSAtMzAwLjUzNjgxKSIvPg0KPC9zdmc+DQo=">
  </a>

**NoCD** is a lightweight and controllable continuous delivery system implemented by Go.

## Preview

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

## Installation Guide

### Docker

1. Create a configuration file (eg `/data/nocd` folder)

   ```shell
   nano /data/nocd/app.ini
   ```

   Refer to the following for the content of the file (`web_listen = 0.0.0.0:8000` configuration do not change)

2. Run NoCD

   ```shell
   docker run -d --name=nocd -p 8000:8000 -v /data/nocd/:/data/conf ghcr.io/naiba/nocd:latest
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

## FAQs

1. Why does my deployment script always fail to execute or not executed at all?

    > Please check whether your PATH path is imported, it is recommended to export the path in advance, it will not be automatically deployed
    >
    > `source .bash_profile`.

2. How to keep running in the background?

    > You can use `systemd`. It is more recommended to run in docker mode.

## License

MIT