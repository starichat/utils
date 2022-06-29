package pool

import (
	"context"
	"testing"
	hello "utils/grpc/proto"
)

func BenchmarkPool_Get(b *testing.B) {
	b.ResetTimer()
	b.SetParallelism(8)
	p := NewConnPool(10, "10.20.43.34:10010")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// do something
			WithPool(p)
		}
	})
	b.StopTimer()
}

func BenchmarkNoPool(b *testing.B) {
	b.ResetTimer()
	b.SetParallelism(8)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// do something
			WithNoPool()
		}
	})
	b.StopTimer()
}

func WithNoPool() {
	conn, _ := DialGrpcConn("10.20.43.34:10010")
	cc := hello.NewGreeterClient(conn.Value())
	_, err := cc.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		return
	}

}

func WithPool(p ConnPool) {
	c, _ := p.Get()
	cc := hello.NewGreeterClient(c.Value())
	_, err := cc.SayHello(context.Background(), &hello.HelloRequest{Name: "astar"})
	if err != nil {
		return
	}
}
