package pool

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
)

//ConnState 连接状态分为
type ConnState int

const (
	Idle       ConnState = iota + 1 //空闲）
	Running                         //运行，正在使用
)

//Conn 一个连接大概有如下功能：
//2. 关闭连接
//3. 连接状态变更
type Conn interface {
	Value() *grpc.ClientConn //对应的类型需要用interface{}来断言
	Close() error//关闭连接
	UpdateState(targetState ConnState) error //连接的状态变更
	Release()  //释放连接
}

type GrpcConnWrap struct {
	*grpc.ClientConn //集成grpc连接
	State ConnState
	mu sync.Mutex
}



//DialGrpcConn 新建一个grpc连接
func DialGrpcConn(addr string) (*GrpcConnWrap, error) {
	cc, err := grpc.DialContext(context.Background(),addr,grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &GrpcConnWrap{ClientConn:cc,State: Idle,mu:sync.Mutex{},
	}, nil
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
	if c.State == targetState {
		//nothing to do
		return ConnStateNotUpdate
	} else if c.State == Idle && targetState == Running {
		c.State = Running

	} else if c.State == Running && targetState == Idle {
		c.State = Idle
	} else {
		return errors.New("state err")
	}
	return nil
}

func (c *GrpcConnWrap) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.State = Idle
}

func (c *GrpcConnWrap) Value() *grpc.ClientConn {
	return c.ClientConn
}

