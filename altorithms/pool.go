package al

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"utils/pool"
)

type Pool interface {
	Get() (pool.Conn, error) //获取连接
	Put(i pool.Conn) error //放置连接
	Remove(i pool.Conn) error //移除连接
	Close()  //关闭连接池
}


func (p *listPool) Get() (pool.Conn, error) {
	//从空闲池中获取数据
	if p.idleConns != nil {
		return p.idleConns.Front().Value.(pool.Conn), nil
	}
	return nil, errors.New("no conn")

}

func (p *listPool) Put(i pool.Conn) error {
	if p.idleConns != nil {
		if _, ok := p.conns[i];ok{
			//该连接已经存在了，直接return
			return errors.New("conn already exist")
		}

		el := p.idleConns.PushFront(i)
		p.conns[i] = el
		return nil
	}
	return errors.New("no pool")
}

func (p *listPool) Close() {
	fmt.Println("close todo ")
}

func (p *listPool) Remove(i pool.Conn) error {
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
	idleConns *list.List
	conns map[pool.Conn]*list.Element
}

func InitPool() Pool {
	p := &listPool{
		idleConns: list.New(),
		conns:     make(map[pool.Conn]*list.Element,0),
	}
	//初始化连接
	for i:=0;i<10;i++{
		wc, err := pool.DialGrpcConn("192.168.3.9:10010")
		if err != nil {
			log.Println("err",err)
			continue
		}
		p.Put(wc)
	}
	return p
}




