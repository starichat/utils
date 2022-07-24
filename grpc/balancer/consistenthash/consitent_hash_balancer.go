package consistenthash

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
)

type BalancerBuilder struct {
	sync.RWMutex
	hashFunc        Func
	keys            slots                       //虚拟节点列表，排序好的虚拟节点，便于通过二分算法快速定位到最近的物理节点
	ring            map[uint64]balancer.SubConn //虚拟节点到物理节点的映射
	nodes           map[Node]struct{}           //物理节点映射，判断当前物理节点是否存在
	NumVirtualNodes int                         // 为每台机器在hash圆环上创建多少个虚拟Node
}

const (
	ConsistentHashBalancer = "ConsistentHashBalancer"
)

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
		csEvltr:   nil,
		state:     connectivity.Ready,
		addrInfos: make(map[string]resolver.Address),
		subConns:  make(map[string]balancer.SubConn),
	}
	//todo 后续可以增加一个连接监控管理组件
	return b

}

func (c *BalancerBuilder) Name() string {
	return ConsistentHashBalancer
}

//实现Balancer接口

type consistentHashBalancer struct {
	cc balancer.ClientConn

	csEvltr   *balancer.ConnectivityStateEvaluator
	state     connectivity.State
	addrInfos map[string]resolver.Address
	subConns  map[string]balancer.SubConn

	scStates map[balancer.SubConn]connectivity.State
	picker   balancer.Picker

	resolverErr error
	connErr     error
}

func (c *consistentHashBalancer) Close() {
}

func (c *consistentHashBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	//定义地址集合
	addrsSet := make(map[string]struct{})
	//遍历ClientConnState的所有地址
	for _, a := range state.ResolverState.Addresses {
		addr := a.Addr
		//将地址更新到consistentHashBalancer的地址列表中
		c.addrInfos[addr] = a
		addrsSet[addr] = struct{}{}
		if sc, ok := c.subConns[addr]; !ok {
			//没有获取到子连接，建立新连接
			newSC, err := c.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{})
			if err != nil {
				log.Printf("Consistent Hash Balancer: failed to create new SubConn: %v", err)
				continue
			}
			//建立新连接的地址映射
			c.subConns[addr] = newSC

		} else {
			//取到了连接，尝试更新当前子连接，没有没问题，则会直接返回，否则仍然会出发一系列连接状态更新的操作
			c.cc.UpdateAddresses(sc, []resolver.Address{a})
		}
	}
	//清理废弃的addr对应的子连接
	for a, sc := range c.subConns {
		if _, ok := addrsSet[a]; !ok {
			c.cc.RemoveSubConn(sc)
			delete(c.subConns, a)
		}
	}
	if len(state.ResolverState.Addresses) == 0 {
		c.ResolverError(fmt.Errorf("可用地址为空"))
		return balancer.ErrBadResolverState
	}
	if c.picker == nil {
		c.picker = NewConsistentHashPicker(c.subConns)
	}
	c.cc.UpdateState(balancer.State{ConnectivityState: c.state, Picker: c.picker})

	return nil
}

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

func (c *consistentHashBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	s := state.ConnectivityState

	switch s {
	case connectivity.Idle:

	case connectivity.Shutdown:
		//重新连接
	case connectivity.TransientFailure:
		//抛错
	}
	//todo 子连接状态变更应该如何处理？

	c.cc.UpdateState(balancer.State{ConnectivityState: c.state, Picker: c.picker})
}
