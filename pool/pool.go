package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

/**
连接池，主要实现了如下功能：
1. 闲置连接释放
2. 连接不够时扩容
3. 连接全满时等待
简单来说需要两个对象：
1. 连接池结构
2. 连接结构
*/

type Pool struct {
	capacity int64
	next     int64
	addr   string
	DialTimeout      time.Duration
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
	ClientPoolSize   int
	lock sync.Mutex
	conns []Conn
}

type ConnPool interface {
	Get() (Conn, error)
	Put(i Conn, idx int64) error
	Close()
}

func NewConnPool(p *Pool) ConnPool {
	return p
}

func (p *Pool) Get() (Conn, error) {
	next := atomic.AddInt64(&p.next, 1)
	idx := next % p.capacity
	//1. 获取索引位置，取出连接
	conn := p.conns[idx]
	//2. 判断当前连接是否可用
	if conn != nil {
		if conn.Status() != UnAvaliable{
			//3. 链接不可用，删除链接
			conn.Close()
		} else {
			return conn, nil
		}
	}

	//4. 创建新的链接,需要加锁
	p.lock.Lock()
	defer p.lock.Unlock()
	cc, err := DialGrpcConn(p.addr)
	if err != nil {
		return nil, err
	}
	p.conns[idx] = cc
	return cc, nil
}

func (p *Pool) Put(c Conn, idx int64) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.conns[idx] = c
	return nil
}

func (p *Pool) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range p.conns {
		if v != nil {
			v.Close()
		}
	}
}


