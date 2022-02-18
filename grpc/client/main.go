package main

import (
	"context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	hello "utils/grpc/proto"
	uB "utils/grpc/balancer"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	//注入自定义负载均衡组件
	//resolver.SetDefaultScheme("dns")
	// 连接服务器
	balancer.Register(base.NewBalancerBuilder("color", uB.MyBalancerBuild{}, base.Config{HealthCheck: false}))
	conn, err := grpc.Dial(":10010", grpc.WithInsecure())
	//conn, err := grpc.DialContext(context.Background(),"",grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
		return
	}
	defer conn.Close()
	//todo, 抓包建立连接

	c := hello.NewGreeterClient(conn)

	// 调用服务端的SayHello
	r, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Println("Greeting", r.Message)
	//再调用服务端程序
	r, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s !\n", r.Message)
}