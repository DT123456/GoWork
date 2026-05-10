package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"

	"gotour/ch22/server"
)

// RPC 客户端示例
// 使用 net/rpc 库调用远程服务
func main() {
	// 连接 RPC 服务器
	// client, err := rpc.Dial("tcp", "localhost:1234")// 基于 TCP 的RPC
	// client, err := rpc.DialHTTP("tcp", "localhost:1234")//新增的，支持HTTP协议的RPC调用
	client, err := jsonrpc.Dial("tcp",  "localhost:1234") // jsonrpc 并没有实现基于 HTTP的传输
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	// 准备参数
	args := server.Args{A: 7, B: 8}

	// 调用远程方法 Add
	var reply int
	err = client.Call("MathService.Add", args, &reply)
	if err != nil {
		log.Fatal("MathService.Add error:", err)
	}

	// 调用远程方法 Multiply
	var replyMul int
	err = client.Call("MathService.Multiply", args, &replyMul)
	if err != nil {
		log.Fatal("MathService.Multiply error:", err)
	}

	fmt.Printf("MathService.Add: %d + %d = %d\n", args.A, args.B, reply)
	fmt.Printf("MathService.Multiply: %d * %d = %d\n", args.A, args.B, replyMul)
}
