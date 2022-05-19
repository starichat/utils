package pool

import (
<<<<<<< HEAD
	"container/ring"
	"context"
=======
	"container/list"
>>>>>>> 11fb87368dcbcb9ed92c7343b0887db5c89c7f29
	"errors"
	"fmt"
	"log"
	"sync"
)

<<<<<<< HEAD
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
=======
type Pool interface {
	Get() (Conn, error) //获取连接
	Put(i Conn) error //放置连接
	Remove(i Conn) error //移除连接
	Close()  //关闭连接池
>>>>>>> 11fb87368dcbcb9ed92c7343b0887db5c89c7f29
}

/**
如果当前连接busy，则将其移到尾部
*/

func (p *listPool)Get() (Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//从空闲池中获取数据
	if p.idleConns != nil  {
		for item := p.idleConns.Front();item != nil;item = item.Next(){
			if cc := item.Value.(Conn); cc.Status() == Running || cc.Status() == Idle {
				state := cc.AddConnStream()
				if state == Busying {
					p.idleConns.MoveToBack(item)
				}
				return cc, nil
			}
			fmt.Println(errors.New("busying"))
			continue

		}
	}
	return nil, errors.New("no conn")

}

func (p *listPool) Put(i Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.idleConns != nil {
		if _, ok := p.conns[i];ok{
			//该连接已经存在了，直接return
			return errors.New("conn already exist")
		}
		fmt.Println("put a conn")
		el := p.idleConns.PushFront(i)
		p.conns[i] = el
		return nil
	}
	return errors.New("no pool")
}

func (p *listPool) Close() {
	fmt.Println("close todo ")
}

func (p *listPool) Remove(i Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if i != nil {
		p.idleConns.Remove(p.conns[i])
		delete(p.conns, i)
	}
	return nil
}

func (p *listPool) AutoResize() error {
	//空闲比率小于一定阈值
	if float32(p.idleConns.Len()) / float32(len(p.conns)) < 0.9 {
		//todo 缩容
	}
	return nil
}



type listPool struct {
	mu *sync.Mutex
	idleConns *list.List
	conns map[Conn]*list.Element
	maxCurrentConn int
}

func InitPool(size int) Pool {
	p := &listPool{
		mu: &sync.Mutex{},
		idleConns: list.New(),
		conns:     make(map[Conn]*list.Element,0),
	}
	//初始化连接
	for i:=0;i<size;i++{
		wc, err := DialGrpcConn("192.168.3.3:10010")
		if err != nil {
			log.Println("err",err)
			continue
		}
		err = p.Put(wc)
		if err != nil {
			panic(err)
		}
	}
	return p
}




