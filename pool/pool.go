package pool
//
//import (
//	"container/ring"
//	"context"
//	"errors"
//	"google.golang.org/grpc"
//	"log"
//	"sync"
//	"time"
//	"container/list"
//)
//
///**
//
//*/
//
//
//type PoolInterface interface {
//	Get() (Conn, error)
//	Close() error
//	Put(c Conn) error
//}
//
//
//type Pool struct {
//	MaxIdle int //最大空闲连接
//	MaxActive int //最大活跃连接
//	MaxCurrentStream int //最大并发stream数量，即单连接下的最大请求数量
//	mu     sync.Mutex    // mu protects the following fields
//	closed bool          // set to true when the pool is closed.
//	active int           // the number of open connections in the pool
//	ch     chan struct{} // limits open connections when p.Wait is true
//	conns     *ring.Ring  //connections
//	current int
//}
//
//type Option func(p *Pool)
//
//var DefaultPool = &Pool{
//	MaxIdle:         64,
//	MaxActive:       64,
//	MaxCurrentStream: 100,
//	mu:              sync.Mutex{},
//	closed:          false,
//	active:          0,
//	ch:              make(chan struct{}),
//	conns: ring.New(64),
//	current: 0,
//}
//
//
//func Dial(address string) (*grpc.ClientConn, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
//	defer cancel()
//	return grpc.DialContext(ctx, address, grpc.WithInsecure())
//}
//
////New...
//func New(address string, options ...Option) (*Pool, error) {
//	if address == "" {
//		panic("invalid address")
//	}
//
//	//构建很多个连接
//	pool := DefaultPool
//	//基于选项模式更新pool结构
//	for _, o := range options {
//		//todo, 选型配置函数待优化
//		o(pool)
//	}
//	//初始化连接
//	for i:=0;i<pool.MaxIdle;i++{
//		cc, err := DialGrpcConn(address, 120 * time.Second)
//		if err !=nil {
//			log.Println("err", err)
//			continue
//		}
//
//		err = pool.Put(cc)
//		if err != nil {
//			log.Println("err", err)
//			continue
//		}
//	}
//	return pool, nil
//}
//
//
//func (p *Pool) Get() (Conn, error) {
//	if p.closed {
//		return nil, errors.New("conn cloesd")
//	}
//	p.mu.Lock()
//	defer p.mu.Unlock()
//	//控制超时.120s内未获取到连接，则返回错误
//	//随机取出当前连接的并发连接数，如果连接还能处理，则随机取出第一个连接，否则取其他空闲连接
//
//}
//
//// Close releases the resources used by the pool.
//func (p *Pool) Close() error {
//	p.closed = true
//	//todo, 释放相关资源
//	return nil
//}
//
//func (p *Pool) Put(c Conn) error {
//	p.mu.Lock()
//	p.mu.Unlock()
//	p.conns.Do()
//	return nil
//}
//
//
//
//
//
