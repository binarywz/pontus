package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (server *Server) Handler(conn net.Conn) {
	// 处理连接
	fmt.Println("connect success...")
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
