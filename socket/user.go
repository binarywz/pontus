package main

import "net"

type User struct {
	Name string
	Addr string
	Chan chan string
	conn net.Conn
}

// NewUser 创建一个用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		Chan: make(chan string),
		conn: conn,
	}

	// 启动监听channel消息的goroutine
	go user.Listen()

	return user
}

// Listen 监听channel消息
func (user *User) Listen() {
	for {
		msg := <-user.Chan
		user.conn.Write([]byte(msg + "\n"))
	}
}
