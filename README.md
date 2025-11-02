# IM-System 即时聊天示例

作者/来源
- 本项目代码参考自 B站 up主 刘丹冰 Aceld 的 Go 聊天室教程视频：
  https://www.bilibili.com/video/BV1gf4y1r79E/?spm_id_from=333.1007.top_right_bar_window_custom_collection.content.click

这是一个用 Go 编写的简单命令行即时聊天（IM）示例项目，将 server 与 client 拆分为两个可独立运行的程序。

目录结构（核心文件）：

- `server/`：服务端程序（包含 `main.go`、`server.go`、`user.go`）
- `client/`：客户端程序（`main.go`）

简介
- 服务端监听 TCP，管理在线用户并负责消息广播与私聊转发。
- 客户端通过交互式终端输入命令和消息，与服务端通信。

快速开始（Windows PowerShell）

1. 在项目根启动服务端（默认监听 127.0.0.1:8888）

```powershell
go run ./server
```

2. 在另一个终端启动客户端（默认连接 127.0.0.1:8888）

```powershell
go run ./client
```

客户端参数
- `-ip`：服务器 IP，默认 `127.0.0.1`
- `-port`：服务器端口，默认 `8888`

例如：

```powershell
go run ./client -ip 127.0.0.1 -port 8888
```

客户端使用说明（交互式）
- 1：公聊模式，向所有在线用户发送消息。
- 2：私聊模式，先查看在线用户（`who`），然后使用 `to|用户名|消息` 格式发送。
- 3：修改用户名，使用 `rename|新名字`。
- 0：退出客户端。

服务端/客户端开发注意
- 代码分为两个独立的 package main，可在同一仓库中分别运行或构建。这样避免了在同一包中存在多个 `main()` 导致冲突。
- 项目已包含 `go.mod` 文件，模块名：`github.com/Su1f4t3/IM-System`。如需改名，请运行 `go mod edit -module <new>` 或 `go mod init <new>`。

调试与构建
- 构建 server 可运行：

```powershell
go build -o im-server ./server
```

- 构建 client 可运行：

```powershell
go build -o im-client ./client
```

常见问题
- 如果提示端口被占用，请确认没有其他进程占用 8888 或修改 server 中的监听地址（`server/main.go`）。
- 若出现连接失败，请检查防火墙或是否用正确的 IP/端口 启动 client。

