package consistenthash

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"log"
)

//BalancerBuilder 负载均衡组件
type BalancerBuilder struct {
}

func init() {
	balancer.Register(newBalancerBuilder())
}

func newBalancerBuilder() *BalancerBuilder {
	return &BalancerBuilder{}
}

//实现Builder接口

func (c *BalancerBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	//构建一个Balancer接口对象
	b := &consistentHashBalancer{
		cc:        cc,
		state:     connectivity.Ready,
		addrInfos: make(map[string]resolver.Address),
		subConns:  make(map[string]balancer.SubConn),
	}
	//todo 后续可以增加一个连接监控管理组件
	return b

}

func (c *BalancerBuilder) Name() string {
	return "ConsistentHashBalancer"
}

//实现Balancer接口

type consistentHashBalancer struct {
	cc balancer.ClientConn //ClientConn

	state     connectivity.State          //连接状态
	addrInfos map[string]resolver.Address //保存resolver解析得到的地址列表
	subConns  map[string]balancer.SubConn //服务端和可用子连接的映射，这里用服务端地址作为key，同时该key也是一致性哈希里的物理node，每个服务端地址维护一份可用连接

	scStates map[balancer.SubConn]connectivity.State //维护子连接及其状态
	picker   balancer.Picker                         //缓存一份picker对象

	resolverErr error
	connErr     error
}

//Close ...
func (c *consistentHashBalancer) Close() {
}

//UpdateClientConnState 当 clientConn 状态变更时会被调用，state 为变更的状态
func (c *consistentHashBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	//定义一个地址集合
	addrs := make(map[string]struct{})
	//遍历ClientConnState中 resolver 提供的最新的所有地址
	for _, a := range state.ResolverState.Addresses {
		addr := a.Addr
		//将地址更新到consistentHashBalancer的地址列表中
		c.addrInfos[addr] = a
		addrs[addr] = struct{}{}
		if sc, ok := c.subConns[addr]; !ok {
			//没有获取到子连接，则基于当前服务端地址，建立一个新连接
			newSC, err := c.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{})
			if err != nil {
				log.Printf("Consistent Hash Balancer: failed to create new SubConn: %v", err)
				continue
			}
			//建立新连接的地址映射
			c.subConns[addr] = newSC

		} else {
			//取到了连接，尝试更新当前子连接，没有没问题，则会直接返回，否则，会尝试保证当前连接的可靠性，诸如，重新建立连接之类的
			c.cc.UpdateAddresses(sc, []resolver.Address{a})
		}
	}
	//清理废弃的addr对应的子连接
	for a, sc := range c.subConns {
		if _, ok := addrs[a]; !ok {
			c.cc.RemoveSubConn(sc)
			delete(c.subConns, a)
		}
	}
	if len(state.ResolverState.Addresses) == 0 {
		c.ResolverError(fmt.Errorf("可用地址为空"))
		return balancer.ErrBadResolverState
	}
	//picker 为空，则建立一个新的picker
	if c.picker == nil {
		c.picker = NewConsistentHashPicker(c.subConns)
	}
	//todo， 判断连接状态
	//这就是连接 balancer 和 picker 的重要过程，基于此函数，会更新ClientConn的连接状态，并且后续通过picker来pick出连接来
	c.cc.UpdateState(balancer.State{ConnectivityState: c.state, Picker: c.picker})

	return nil
}

// ResolverError 当reslover组件报告某些错误，可通过该函数回调
func (c *consistentHashBalancer) ResolverError(err error) {
	c.resolverErr = err

	//todo 处理picker
	if c.state != connectivity.TransientFailure {
		// The picker will not change since the balancer does not currently
		// report an error.
		return
	}
	c.cc.UpdateState(balancer.State{
		ConnectivityState: c.state,
		Picker:            c.picker,
	})
}

// UpdateSubConnState 当子连接状态变更会回调该函数
func (c *consistentHashBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	s := state.ConnectivityState

	switch s {
	case connectivity.Idle:

	case connectivity.Shutdown:
		//重新连接
		sc.Connect()
	case connectivity.TransientFailure:
		//抛错
	}

	c.cc.UpdateState(balancer.State{ConnectivityState: c.state, Picker: c.picker})
}
