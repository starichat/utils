package balancer

import (
	"google.golang.org/grpc/balancer"
)

/**
以下完成一套基于hash算法的负载均衡策略
*/




type MyBalancer struct {
	cc balancer.ClientConn
}

type MyBalancerBuild struct {
}

func (*MyBalancerBuild) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &MyBalancer{cc: cc}
}

func (*MyBalancerBuild) Name() string {
	return "mybalancer"
}

func (mp *MyBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	panic("implement me")
}

func (mp *MyBalancer) ResolverError(err error) {
	panic("implement me")
}

func (mp *MyBalancer) UpdateSubConnState(conn balancer.SubConn, state balancer.SubConnState) {
	panic("implement me")
}

func (mp *MyBalancer) Close() {
	panic("implement me")
}

