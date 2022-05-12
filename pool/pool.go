package pool

import (
	"container/ring"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

/**
连接池，主要实现了如下功能：
1. 连接的自动扩容和缩容算法
2. 连接的调度
简单来说需要两个对象：
1. 连接池结构
2. 连接结构
*/



type Pool struct {

	MaxIdle int //最大空闲连接
	MaxActive int //最大活跃连接
	IdleTimeout time.Duration //todo
	Wait bool //如果当前连接池，没有空闲连接了，就会可用，整个池子就进入等待状态了。这里可以用chan来处理，还可以规避并发问题

	MaxConnLifetime time.Duration //连接最长生命周期

	mu     sync.Mutex    // mu protects the following fields
	closed bool          // set to true when the pool is closed.
	active int           // the number of open connections in the pool
	ch     chan struct{} // limits open connections when p.Wait is true
	conns  *ring.Ring     // idle connections
	idle *ring.Ring //idle 指向同一个环，只是索引位置不同

}

type Option func(Pool *Pool)

var DefaultPool = &Pool{
	MaxIdle:         10,
	MaxActive:       10,
	IdleTimeout:     0,
	Wait:            false,
	MaxConnLifetime: 0,
	mu:              sync.Mutex{},
	closed:          false,
	active:          0,
	ch:              make(chan struct{}),
	idle:            []Conn,
}


func Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithInsecure())
}

//New,toodo 改造为options
func New(address string, options Option...) (*Pool, error) {
	if address == "" {
		panic("invalid address")
	}

	//构建很多个连接
	pool := DefaultPool
	//基于选项模式更新pool结构
	for _, o := range options {
		o(pool)
	}
	//基于options初始化配置参数
	for i:=0;i<pool.MaxActive;i++{
		cc, err := DialGrpcConn(address)
		if err !=nil {
			log.Println("err", err)
			continue
		}

		err = pool.Put(cc)
		if err != nil {
			log.Println("err", err)
			continue
		}
	}
	return pool, nil
}


func (p *Pool) Get() (Conn, error) {
	if p.closed {
		return nil, errors.New("conn cloesd")
	}
	if p.Wait  {
		return nil, errors.New("conn busying")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if  > 0 {
		cc :=<- p.idle
		return cc, nil
	}
	return nil, errors.New("conn busying")
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	return nil
}

func (p *Pool) Put(c Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idle <- c
	return nil
}




// PoolStats contains pool statistics.
type PoolStats struct {
	// ActiveCount is the number of connections in the pool. The count includes
	// idle connections and connections in use.
	ActiveCount int
	// IdleCount is the number of idle connections in the pool.
	IdleCount int
}

// Stats returns pool's statistics.
func (p *Pool) Stats() PoolStats {
	p.mu.Lock()
	stats := PoolStats{
		ActiveCount: p.active,
		IdleCount:   len(p.idle),
	}
	p.mu.Unlock()

	return stats
}

// ActiveCount returns the number of connections in the pool. The count
// includes idle connections and connections in use.
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

// IdleCount returns the number of idle connections in the pool.
func (p *Pool) IdleCount() int {
	p.mu.Lock()
	idle := len(p.idle)
	p.mu.Unlock()
	return idle
}




