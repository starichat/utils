package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"time"
	al "utils/altorithms"
	hello "utils/grpc/proto"
	"utils/pool"
)

func main() {
	//注入自定义负载均衡组件
	//for {
	//	if time.Now().Minute() == 03 {
	//		break
	//	}
	//}


	//resolver.SetDefaultScheme("dns")
	// 连接服务器
	//balancer.Register(&uB.MyBalancerBuild{})
	//conn, err := grpc.Dial(":10010", grpc. WithDefaultServiceConfig(`{"loadBalancingPolicy":"color"}`))
	ppp := al.InitPool(100)
	//conn, err := grpc.DialContext(context.Background(),"192.168.3.3:10010",grpc.WithInsecure())
	//if err != nil {
	//	fmt.Printf("faild to connect: %v", err)
	//	return
	//}
	//cc, err := grpc.DialContext(context.Background(),"192.168.3.3:10010",grpc.WithInsecure())
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("start")
	t := time.Now().Unix()
	//todo, 抓包建立连接
	//conn 使用连接池来构建
	wg := sync.WaitGroup{}
	wg.Add(100000)
	for i:=0;i<100000;i++{
		go func() {
			//var conn pool.Conn
			//var err error
			//for {
			//	conn, err  = ppp.Get()
			//	if err != nil {
			//		fmt.Printf("could not get conn: %v\n", err)
			//		continue
			//	}
			//	break
			//}
			//
			//
			//c := hello.NewGreeterClient(conn.Value())
			//t := time.Now().UnixMilli()
			//_, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
			//if err != nil {
			//	fmt.Printf("could not dail greet: %v", err)
			//	conn.Release()
			//	wg.Done()
			//	return
			//}
			//
			//fmt.Printf("main1 Greeting success and consumer %d\n", time.Now().UnixMilli()-t)
			//conn.Release()
			//wg.Done()
			//withSingleConn(&wg,cc)
			withPool(&wg,ppp)
		}()

	}
	wg.Wait()
	fmt.Println("耗时",time.Now().Unix()-t)
	//
	//// 调用服务端的SayHello
	//
	////再调用服务端程序
	//r, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	//if err != nil {
	//	fmt.Printf("could not greet: %v", err)
	//}
	//fmt.Printf("Greeting: %s !\n", r.Message)
}

func withPool(wg *sync.WaitGroup, p al.Pool) {
	t := time.Now().UnixMilli()
	var conn pool.Conn
	var err error
	for {
		conn, err  = p.Get()
		if err != nil {
			fmt.Printf("could not get conn: %v\n", err)
			continue
		}
		break
	}


	c := hello.NewGreeterClient(conn.Value())

	_, err = c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not dail greet: %v", err)
		conn.Release()
		wg.Done()
		return
	}

	fmt.Printf("main1 Greeting success and consumer %d\n", time.Now().UnixMilli()-t)
	conn.Release()
	wg.Done()
}

func withSingleConn(wg *sync.WaitGroup,cc *grpc.ClientConn) {

	c := hello.NewGreeterClient(cc)
	t := time.Now().UnixMilli()
	_, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		fmt.Printf("could not dail greet: %v", err)
		wg.Done()
		return
	}
	fmt.Printf("main1 Greeting success and consumer %d\n", time.Now().UnixMilli()-t)
	wg.Done()
}