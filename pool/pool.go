package pool

import "sync"

/**
连接池，主要实现了如下功能：
1. 连接的自动扩容和缩容算法
2. 连接的调度
简单来说需要两个对象：
1. 连接池结构
2. 连接结构
 */

type Pool struct {
	mu      sync.Mutex
	minConn int // 最小连接数
	maxConn int // 最大连接数
	numConn int // 池已申请的连接数
	conns   chan *Conn //当前池中空闲连接实例(用chan这种go语言特性的管道来做连接池结构非常优雅，可以很好地控制并发，也能支持容量的扩展)
	close   bool
	scheduler Schedule
}

//Pooler 这里需要参考一下现有的主流连接池的设计
type Pooler interface {
	Close()
	Get()
	Put()
}

