package pool

import (
	"context"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

//ConnState 连接状态分为
type ConnState int

const (
	Idle       ConnState = iota + 1 //空闲
	Running                         //运行，正在使用
	Busying                         //忙碌，暂不可用
)

//Conn 一个连接大概有如下功能：
//2. 关闭连接
//3. 连接状态变更
type Conn interface {
	Value() *grpc.ClientConn //对应的类型需要用interface{}来断言
	Close() error//关闭连接
	UpdateState(targetState ConnState) error //连接的状态变更
	Release() ConnState //释放连接
	AddConnStream() ConnState //增加连接
	Status() ConnState //当前连接状态
}

type GrpcConnWrap struct {
	*grpc.ClientConn //集成grpc连接
	State ConnState
	mu sync.Mutex
	ConnStreamCount int64 //当前连接数量
	maxSize int64
}

func (c *GrpcConnWrap) Status() ConnState {
	return c.State
}



//DialGrpcConn 新建一个grpc连接
func DialGrpcConn(addr string) (*GrpcConnWrap, error) {
	cc, err := grpc.DialContext(context.Background(),addr,grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &GrpcConnWrap{ClientConn:cc,State: Idle,mu:sync.Mutex{},
	maxSize: 10000}, nil
}



//Close 关闭连接
func (c *GrpcConnWrap) Close() error {
	return c.ClientConn.Close()
}



//UpdateState 更新连接状态,连接更新操作需要加锁
//ready -> idle
//idle -> ready
//idle -> busying
func (c *GrpcConnWrap)	UpdateState(targetState ConnState) error {//连接的状态变更
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ConnStreamCount < c.maxSize {

	}

	return nil
}

func (c *GrpcConnWrap) Release() ConnState {
	atomic.AddInt64(&c.ConnStreamCount,-1)
	c.State = Running

	return c.State
}

func (c *GrpcConnWrap) Value() *grpc.ClientConn {
	return c.ClientConn
}

func (c *GrpcConnWrap) AddConnStream() ConnState {

	atomic.AddInt64(&c.ConnStreamCount,1)
	if c.ConnStreamCount < c.maxSize {
		c.State = Running
	} else {
		c.State = Busying
	}
	return c.State
}
