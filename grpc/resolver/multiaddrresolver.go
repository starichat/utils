package resolver

import (
	"context"
	"errors"
	"google.golang.org/grpc/resolver"
	"log"
)

/**
自定义resolver
*/

var Nodes = map[string][]string{
	"grpc-server": []string{
		"127.0.0.1:10000",
		"127.0.0.1:10001",
	},
}

type myResolverBuilder struct {
}

func init() {
	resolver.Register(&myResolverBuilder{})
}

func (m *myResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	r := &myResolver{
		ctx:    nil,
		cancel: nil,
		cc:     cc,
		ch:     make(chan struct{}),
	}
	go r.watch(target.Endpoint)
	go func() {
		//启动前先更新一波resolver地址
		r.ResolveNow(resolver.ResolveNowOptions{})
	}()
	r.getNodes(target.Endpoint)
	return r, nil
}

func (m *myResolverBuilder) Name() string {
	return "myResolver"
}

type myResolver struct {
	addrs  []string
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	ch     chan struct{} //监控地址变化的通道
}

func (m *myResolver) ResolveNow(options resolver.ResolveNowOptions) {
	//启动前先更新一波resolver地址
	select {
	case m.ch <- struct{}{}:
		log.Println("发送变更通知")
	default:
	}

}

func (m *myResolver) Close() {

}

func (m *myResolverBuilder) Scheme() string {
	return "myResolver"
}

func (m *myResolver) watch(name string) {
	//监控通道变化
	for {
		select {
		case <-m.ch: //能够从通道中取出ch，则证明地址有变化
			//触发重新解析地址
			log.Println("获取更新地址指令")
			nodes, err := m.getNodes(name)
			if err != nil {
				return
			}
			log.Println("更新后的node为", nodes.Addresses)
		}
	}

}

func (m *myResolver) getNodes(name string) (resolver.State, error) {
	s := resolver.State{}
	//获取最新地址
	if v, ok := Nodes[name]; ok {
		for _, item := range v {
			//添加地址到resolver中
			s.Addresses = append(s.Addresses, resolver.Address{
				Addr:       item,
				ServerName: name,
				Attributes: nil,
			})
		}
	} else {
		return s, errors.New("no nodes available")
	}
	m.addrs = Nodes[name]
	//更新resolve状态，传递给balancer
	err := m.cc.UpdateState(s)
	if err != nil {
		return resolver.State{}, err
	}
	return s, nil
}
