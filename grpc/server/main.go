package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
<<<<<<< HEAD
	"sync/atomic"
	"time"
=======
>>>>>>> 11fb87368dcbcb9ed92c7343b0887db5c89c7f29
	hello "utils/grpc/proto"
)

type HelloServer struct {

}

func (s *HelloServer) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	fmt.Println("get ",in.GetName())
<<<<<<< HEAD
	time.Sleep(1*time.Second)
	atomic.AddInt64(&i,1)
	fmt.Printf("第%d次请求\n",i)
=======
>>>>>>> 11fb87368dcbcb9ed92c7343b0887db5c89c7f29
	return &hello.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

var i int64

func main() {
	lis, err := net.Listen("tcp", "192.168.3.9:10010")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer() // 创建gRPC服务器
	hello.RegisterGreeterServer(s, &HelloServer{}) // 在gRPC服务端注册服务

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	fmt.Println("start server")
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}