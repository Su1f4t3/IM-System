# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个用 Go 语言实现的简单即时通讯（IM）系统，采用 TCP 协议进行通信。项目包含三个核心文件，实现了基本的聊天室功能。

## 核心架构

### 主要组件

- **Server** (`server.go`): 核心服务器，管理用户连接和消息广播
- **User** (`user.go`): 用户实体，处理单个客户端连接
- **Main** (`main.go`): 程序入口点

### 架构模式

1. **并发处理**: 每个客户端连接都有独立的 goroutine 处理
2. **消息广播**: 使用中央 channel (`Message`) 进行消息分发
3. **用户管理**: 使用 `OnlineMap` 维护在线用户列表，配合读写锁保证并发安全
4. **消息流向**: 客户端 → Server.Handler → Server.Message → User.C → 客户端

### 关键方法流程

- **服务器启动**: `server.Start()` → 启动 TCP 监听器 → 启动 `ListenMessager()` goroutine → 接受客户端连接
- **用户连接**: `server.Handler()` → 创建 `User` 实例 → 调用 `user.Online()` → 启动消息接收 goroutine
- **消息处理**: `user.DoMessage()` → `server.BroadCast()` → 发送到 `server.Message` channel → `ListenMessager()` 分发给所有用户
- **用户下线**: `user.Offline()` → 从 `OnlineMap` 移除 → 广播下线消息

### 关键数据结构

- `Server.OnlineMap`: 存储在线用户的映射表
- `Server.Message`: 消息广播 channel
- `User.C`: 用户专用消息 channel

## 常用命令

### 运行服务器
```bash
go run *.go
```

### 编译项目
```bash
go build -o im-server
```

### 运行编译后的服务器
```bash
./im-server
```

### 格式化代码
```bash
go fmt *.go
```

### 检查代码问题
```bash
go vet *.go
```

## 开发注意事项

- 服务器默认监听地址：127.0.0.1:8888
- 用户名默认使用客户端连接地址
- 消息格式：`[用户地址]用户名: 消息内容`
- 系统自动发送用户上线/下线通知
- 使用简单的 TCP 连接，可使用 telnet 或 netcat 测试

### 特殊命令

- **who**: 查询当前在线用户列表，发送给当前用户
- 其他消息: 广播给所有在线用户

### 错误处理

- 连接断开时自动清理用户并广播下线消息
- 网络读取错误会打印到控制台并关闭连接
- 服务器启动失败会直接退出

## 测试连接

可以使用以下工具连接服务器测试：
```bash
telnet 127.0.0.1 8888
# 或
nc 127.0.0.1 8888
```