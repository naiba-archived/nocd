# NoCD 持续交付系统

![构建状态](https://github.com/naiba/nocd/workflows/Build%20Docker%20Image/badge.svg) <a href="README_en-US.md">
    <img height="20px" src="https://img.shields.io/badge/EN-flag.svg?color=555555&style=flat&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB2aWV3Qm94PSIwIDAgNjAgMzAiIGhlaWdodD0iNjAwIj4NCjxkZWZzPg0KPGNsaXBQYXRoIGlkPSJ0Ij4NCjxwYXRoIGQ9Im0zMCwxNWgzMHYxNXp2MTVoLTMwemgtMzB2LTE1enYtMTVoMzB6Ii8+DQo8L2NsaXBQYXRoPg0KPC9kZWZzPg0KPHBhdGggZmlsbD0iIzAwMjQ3ZCIgZD0ibTAsMHYzMGg2MHYtMzB6Ii8+DQo8cGF0aCBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iNiIgZD0ibTAsMGw2MCwzMG0wLTMwbC02MCwzMCIvPg0KPHBhdGggc3Ryb2tlPSIjY2YxNDJiIiBzdHJva2Utd2lkdGg9IjQiIGQ9Im0wLDBsNjAsMzBtMC0zMGwtNjAsMzAiIGNsaXAtcGF0aD0idXJsKCN0KSIvPg0KPHBhdGggc3Ryb2tlPSIjZmZmIiBzdHJva2Utd2lkdGg9IjEwIiBkPSJtMzAsMHYzMG0tMzAtMTVoNjAiLz4NCjxwYXRoIHN0cm9rZT0iI2NmMTQyYiIgc3Ryb2tlLXdpZHRoPSI2IiBkPSJtMzAsMHYzMG0tMzAtMTVoNjAiLz4NCjwvc3ZnPg0K">
</a>

**NoCD** 是一个轻量可控的持续交付系统。

## 界面预览

| ![首页截图](https://github.com/naiba/nocd/raw/master/README/首页截图.png) | ![服务器管理](https://github.com/naiba/nocd/raw/master/README/服务器管理.png) | ![项目管理](https://github.com/naiba/nocd/raw/master/README/项目管理.png) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![交付记录](https://github.com/naiba/nocd/raw/master/README/交付记录.png) | ![管理中心](https://github.com/naiba/nocd/raw/master/README/查看日志.png) | ![查看日志](https://github.com/naiba/nocd/raw/master/README/管理中心.png)  |

## 功能特色

- 服务器：可以添加多个部署服务器
- 项目：支持解析各种流行 Git 托管平台的 Webhook
- 通知：灵活的自定义 Webhook
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

   ```shell
   docker run -d --name=nocd -p 8000:8000 -v /data/nocd/:/data/conf ghcr.io/naiba/nocd:latest
   ```

### 源码编译

1. Clone 源代码

2. 进入应用目录 `cd nocd/cmd/web`

3. 编译二进制

   ```shell
   go build
   ```

4. 在 `conf/app.ini` 创建配置文件

   ```ini
   [nocd]
   cookie_key_pair = i_love_NoCD
   debug = true
   domain = your_domain_name # or ip:port
   web_listen = 0.0.0.0:8000
   loc = Asia/Shanghai
   [third_party]
   google_analysis = "NB-XXXXXX-1" # optional
   github_oauth2_client_id = example
   github_oauth2_client_secret = example
   sentry_dsn = "https://example:xx@example.io/project_id" # optional
   ```

5. 运行

   ```shell
   ./web
   ```

6. 在 `GitHub` 设置回调：`http(s)://your_domain_name/oauth2/callback`

## 常见问题

1. 为什么我的部署脚本总是执行失败 或者 根本没有执行？

    > 请检查您的 PATH 路径是否引入，建议提前 export 一下路径，自动部署的时候不会
    >
    > `source .bash_profile`。

2. 如何保持后台运行？

    > 可以使用 `systemd` 。 更推荐使用docker方式运行。

## License

MIT
