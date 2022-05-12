package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
	hello "utils/grpc/proto"
)

type HelloServer struct {

}

func (s *HelloServer) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	fmt.Println("get ",in.GetName())
	time.Sleep(1*time.Second)
	return &hello.HelloReply{Message: "Hello again " + in.GetName()}, nil
}


func main() {
	lis, err := net.Listen("tcp", ":10010")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer() // 创建gRPC服务器
	hello.RegisterGreeterServer(s, &HelloServer{}) // 在gRPC服务端注册服务

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}