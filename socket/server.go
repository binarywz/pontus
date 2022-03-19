package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// Listen 监听Message广播消息channel的goroutine
func (server *Server) Listen() {
	for {
		msg := <-server.Message

		// 将msg发送给所有的在线User
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.Chan <- msg
		}
		server.mapLock.Unlock()
	}
}

// BroadCast 广播消息
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// 处理连接
	userAddr := conn.RemoteAddr().String()
	user := NewUser(conn)
	fmt.Println(userAddr, "online...")

	// 用户上线
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线消息
	server.BroadCast(user, "online")

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, readErr := conn.Read(buf)
			if n == 0 {
				server.BroadCast(user, "offline")
				return
			}
			if readErr != nil {
				fmt.Println("Conn read err:", readErr)
				return
			}

			// 提取用户的消息
			msg := string(buf)

			// 广播消息
			server.BroadCast(user, msg)
		}
	}()

	// 当前handler阻塞
	select {}
}

// Start 启动服务的接口
func (server *Server) Start() {
	fmt.Println("server start...")
	// socket listen
	listener, listenErr := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if listenErr != nil {
		fmt.Println("net listen err: ", listenErr)
		return
	}
	// close socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go server.Listen()

	for {
		// accept
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			fmt.Println("listen accept err: ", acceptErr)
			continue
		}

		// do handler
		go server.Handler(conn)
	}

}
