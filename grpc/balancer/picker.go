package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
	"sync"
	"sync/atomic"
)

type MyPicker struct {
	lock  sync.Locker
	conns []*subConn
}

//subConn 子连接
type subConn struct {
	addr     resolver.Address
	conn     balancer.SubConn
	lag      uint64 // 用来保存 ewma 值
	inflight int64  // 用在保存当前节点正在处理的请求总数
	success  uint64 // 用来标识一段时间内此连接的健康状态
	requests int64  // 用来保存请求总数
	last     int64  // 用来保存上一次请求耗时, 用于计算 ewma 值
	pick     int64  // 保存上一次被选中的时间点
}

func (mp *MyPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {

	mp.lock.Lock()
	defer mp.lock.Unlock()
	var chosen *subConn
	switch len(mp.conns) {
	//case 0:
	//	return nil, nil, balancer.ErrNoSubConnAvailable // 没有可用链接
	//case 1:
	//	chosen = mp.choose(p.conns[0], nil) // 只有一个链接
	//case 2:
	//	chosen = p.choose(p.conns[0], p.conns[1])
	//default: // 选择一个健康的节点
	//	var node1, node2 *subConn
	//	for i := 0; i < pickTimes; i++ {
	//		a := p.r.Intn(len(p.conns))
	//		b := p.r.Intn(len(p.conns) - 1)
	//		if b >= a {
	//			b++
	//		}
	//		node1 = p.conns[a]
	//		node2 = p.conns[b]
	//		if node1.healthy() && node2.healthy() {
	//			break
	//		}
	//	}
	//	chosen = p.choose(node1, node2)
	//todo，具体的负载策略，从conns中选取出一个连接
	}
	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)
	return balancer.PickResult{
		SubConn: mp.conns[0],
		Done:    nil,
	}, nil
}

func (p *MyPicker) buildDoneFunc(c *subConn) func(info balancer.DoneInfo) {

	return func(info balancer.DoneInfo) {

	}
}

func (s subConn) UpdateAddresses(addresses []resolver.Address) {
	panic("implement me")
}

func (s subConn) Connect() {
	panic("implement me")
}
