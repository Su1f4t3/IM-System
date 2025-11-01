package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// 定义服务器结构体
type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User // 在线用户列表
	maplock   sync.RWMutex     // 读写锁，保护并发访问
	Message   chan string      // 消息广播channel
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听消息（Message）广播channel的goroutine，一旦有消息，就将消息发送给所有在线用户（User）的channel
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		// 将msg发送给所有在线的User
		s.maplock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.maplock.Unlock()
	}
}

// 广播消息方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg

	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {

	// 创建新用户对象
	user := NewUser(conn, s)

	// 用户上线业务
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096) // 创建一个 4096 字节的缓冲区（切片）用于存储从客户端接收的数据

		// 无限循环等待客户端发送消息
		for {
			n, err := conn.Read(buf) // 会阻塞直到接收到数据或连接断开

			// 如果 n 等于 0，表示客户端断开连接
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1]) // 去掉换行符 "\n"

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户是活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，什么都不做，继续阻塞

		case <-time.After(time.Minute * 30):
			// 已经超时
			user.SendMsg("你被踢下线了，因为你长时间没有活动。。。\n")

			// 销毁用户资源
			user.Offline()
			close(user.C)

			// 关闭连接
			conn.Close()

			// 退出当前handler
			return
		}
	}
}

// 启动服务器的接口
func (s *Server) Start() {
	// 创建一个 TCP 网络监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	// 关闭监听器
	defer listener.Close()

	// 启动监听Message的goroutine
	go s.ListenMessager()

	// 无限循环等待客户端连接
	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue // 继续等待下一个连接
		}

		// 为每个连接启动独立的处理goroutine
		go s.Handler(conn)
	}
}
