package pool

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

//ConnState 连接状态分为
type ConnState int

const (
	UnAvaliable = 0
	Avaliable
)

//Conn 一个连接大概有如下功能：
//2. 关闭连接
//3. 连接状态变更
type Conn interface {
	Value() *grpc.ClientConn //对应的类型需要用interface{}来断言
	Close() error            //关闭连接
	Status() ConnState       //当前连接状态
}

type GrpcConnWrap struct {
	*grpc.ClientConn //集成grpc连接
	State            ConnState
}

func (c *GrpcConnWrap) Status() ConnState {
	switch c.GetState() {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return UnAvaliable
	default:
		return Avaliable
	}
	return c.State
}

//DialGrpcConn 新建一个grpc连接
func DialGrpcConn(addr string) (*GrpcConnWrap, error) {
	cc, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &GrpcConnWrap{ClientConn: cc, State: Avaliable}, nil
}

//Close 关闭连接
func (c *GrpcConnWrap) Close() error {
	return c.ClientConn.Close()
}

func (c *GrpcConnWrap) Value() *grpc.ClientConn {
	return c.ClientConn
}
