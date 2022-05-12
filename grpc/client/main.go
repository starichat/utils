package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	hello "utils/grpc/proto"
)

func main() {
	//注入自定义负载均衡组件
	//resolver.SetDefaultScheme("dns")
	// 连接服务器
	//balancer.Register(&uB.MyBalancerBuild{})
	//conn, err := grpc.Dial("198.168.3.9:10010", grpc. WithDefaultServiceConfig(`{"loadBalancingPolicy":"color"}`))
	conn, err := grpc.DialContext(context.Background(),"198.168.3.9:10010",grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
		return
	}
	c := hello.NewGreeterClient(conn)
	// 调用服务端的SayHello
	r, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Println(r)
	defer conn.Close()
	//todo, 抓包建立连接
	//conn 使用连接池来构建
	for i:=0;i<1;i++{

		//fmt.Println("Greeting", r.Message)
		////再调用服务端程序
		//r, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
		//if err != nil {
		//	fmt.Printf("could not greet: %v", err)
		//}
		//fmt.Printf("Greeting: %s !\n", r.Message)
	}


}