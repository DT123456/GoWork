package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"gotour/ch22/server"
)

// RPC 服务器示例
// 注册 MathService 并监听客户端请求
func main() {
	// 注册 RPC 服务
	err := rpc.RegisterName("MathService", new(server.MathService))
	if err != nil {
		log.Fatal("register error:", err)
	}

	// 处理 HTTP 请求（可选，支持基于 HTTP 的 RPC）
	// rpc.HandleHTTP()//新增的

	// 监听 TCP 端口
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	// http.Serve(l, nil)//换成http的服务器，支持HTTP协议的RPC调用
	
	defer l.Close()

	log.Println("RPC server listening on :1234")

	// 接受客户端连接（阻塞调用，在 goroutine 中运行）
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("jsonrpc.Serve: accept:", err.Error())
			return
		}
		//json rpc
		go jsonrpc.ServeConn(conn)
	}


	// 阻塞主线程，避免程序退出
	select {}
}
