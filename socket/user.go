package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	Chan   chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		Chan:   make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听channel消息的goroutine
	go user.Listen()

	return user
}

// Online 用户上线
func (user *User) Online() {
	// 用户上线
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "online")
}

// Offline 用户下线
func (user *User) Offline() {
	// 用户下线
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "offline")
}

func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		fmt.Println("user DoMessage:", msg)
		// 查询当前在线用户
		user.server.mapLock.Lock()
		for _, item := range user.server.OnlineMap {
			onlineInfo := "[" + item.Addr + "]" + item.Name + ": online..."
			user.SendMsg(onlineInfo)
		}
		user.server.mapLock.Unlock()
	} else {
		user.server.BroadCast(user, msg)
	}

}

// Listen 监听channel消息
func (user *User) Listen() {
	for {
		msg := <-user.Chan
		user.conn.Write([]byte(msg + "\n"))
	}
}
