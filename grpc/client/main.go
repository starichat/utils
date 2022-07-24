package main

import (
	"context"
	"fmt"
	hello "utils/grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	grpcCallWithConsistentBalancer()

	//fmt.Println("start")
	////注入自定义负载均衡组件
	//t := time.Now().Unix()
	////resolver.SetDefaultScheme("dns")
	//// 连接服务器
	////balancer.Register(&uB.MyBalancerBuild{})
	////conn, err := grpc.Dial(":10010", grpc. WithDefaultServiceConfig(`{"loadBalancingPolicy":"color"}`))
	//conn, err := grpc.DialContext(context.Background(), "10.20.43.34:10010", grpc.WithInsecure())
	//if err != nil {
	//	fmt.Printf("faild to connect: %v", err)
	//	return
	//}
	//c := hello.NewGreeterClient(conn)
	//defer conn.Close()
	////todo, 抓包建立连接
	////conn 使用连接池来构建
	//wg := sync.WaitGroup{}
	//wg.Add(50000)
	//for i := 0; i < 50000; i++ {
	//	go func() {
	//		t := time.Now().UnixMilli()
	//		_, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	//		if err != nil {
	//			fmt.Printf("could not greet: %v", err)
	//		}
	//		fmt.Printf("main2 Greeting success and consumer %d\n", time.Now().UnixMilli()-t)
	//		wg.Done()
	//	}()
	//
	//}
	//wg.Wait()
	//fmt.Println("耗时", time.Now().Unix()-t)
	////
	////// 调用服务端的SayHello
	////
	//////再调用服务端程序
	////r, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	////if err != nil {
	////	fmt.Printf("could not greet: %v", err)
	////}
	////fmt.Printf("Greeting: %s !\n", r.Message)
}

func grpcCallWithConsistentBalancer() {
	conn, err := grpc.Dial(":10010", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"ConsistentHashBalancer"}`))
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
		return
	}
	c := hello.NewGreeterClient(conn)
	_, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
}
