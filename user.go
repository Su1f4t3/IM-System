package main

import "net"

type User struct {
	Name string      // 用户名（默认使用连接地址）
	Addr string      // 用户地址
	C    chan string // 接收消息的channel
	conn net.Conn    // 网络连接对象

	server *Server // 当前用户所属的服务器
}

// 创建一个用户的接口
func NewUser(conn net.Conn, server *Server) *User {
	// 获取客户端连接地址作为用户名和地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (u *User) Online() {

	// 用户上线，将用户加入在线用户列表（OnlineMap）中
	u.server.maplock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.maplock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// 用户下线业务
func (u *User) Offline() {

	// 用户下线，将用户从在线用户列表（OnlineMap）中移除
	u.server.maplock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.maplock.Unlock()

	// 广播当前用户下线消息
	u.server.BroadCast(u, "已下线")

	// 统一关闭用户的channel
	close(u.C)
}

// 给当前user对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户处理消息业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {

		// 查询当前在线用户都有哪些
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.maplock.Unlock()

	} else if len(msg) > 7 && msg[0:7] == "rename|" {
		// 消息格式：rename|张三
		newName := msg[7:]

		// 判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名已被使用\n")
		} else {
			u.server.maplock.Lock()
			delete(u.server.OnlineMap, u.Name) // 删除原来的用户名
			u.server.OnlineMap[newName] = u    // 添加新用户名
			u.server.maplock.Unlock()

			u.Name = newName
			u.SendMsg("用户名更新成功，新的用户名为：" + u.Name + "\n")
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}

// 监听用户（User）消息 channel ，一旦有消息就发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
