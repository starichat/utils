# grpc 负载均衡

## 工作原理
grpc balancer 工作原理为：
grpc resolver -> build -> balancer build -> balancer pick

基于此，我们可以轻松构建自己的grpc负载均衡处理组件

常见的负载均衡算法有：
1. 一致性hash
2. p2c

参考 [一致性哈希算法](./consistenthash/README.md) 了解一致性哈希算法的 go 语言实践。

基于上述算法，我们加以丰富，让其成为grpc链接的负载均衡组件。

## 实现过程
上一节从上层到底层讲述了 grpc balancer 的工作原理，这一节便从底层向上层介绍实现过程。
首先需要解决一个问题，就是我们日常使用 grpc 连接去进行 RPC 调用时，只是拿到了一个 grpc 的 ClientConn 接口，相关接口实现方也没有暴露出底层 ConnID 来唯一标识该连接，也就暂时无法利用一致性哈希算法来准确定位到对应的 grpc 连接上去。

再仔细阅读源码，发现一个完成的 grpc 请求，最终是通过 balancer.SubConn 来完成的，是不是可以通过 SubConn 来唯一标识子连接，用来实现 grpc 的负载均衡。
我们要继承实现的也是这一块的代码。

那么 balancer.SubConn 就相当于是哈希算法中的物理节点了。
核心问题解决了，直接套上一致性哈希算法的外衣，直接一把嗦，得到如下结构体
```go


```

1. newBalancerBuild
```go
// Builder creates a balancer.
type Builder interface {
	// Build creates a new balancer with the ClientConn.
	Build(cc ClientConn, opts BuildOptions) Balancer
	// Name returns the name of balancers built by this builder.
	// It will be used to pick balancers (for example in service config).
	Name() string
}
```
2. 基于Builder->Balancer
```go
// Balancer takes input from gRPC, manages SubConns, and collects and aggregates
// the connectivity states.
//
// It also generates and updates the Picker used by gRPC to pick SubConns for RPCs.
//
// HandleSubConnectionStateChange, HandleResolvedAddrs and Close are guaranteed
// to be called synchronously from the same goroutine.
// There's no guarantee on picker.Pick, it may be called anytime.
type Balancer interface {
	// HandleSubConnStateChange is called by gRPC when the connectivity state
	// of sc has changed.
	// Balancer is expected to aggregate all the state of SubConn and report
	// that back to gRPC.
	// Balancer should also generate and update Pickers when its internal state has
	// been changed by the new state.
	//
	// Deprecated: if V2Balancer is implemented by the Balancer,
	// UpdateSubConnState will be called instead.
	HandleSubConnStateChange(sc SubConn, state connectivity.State)
	// HandleResolvedAddrs is called by gRPC to send updated resolved addresses to
	// balancers.
	// Balancer can create new SubConn or remove SubConn with the addresses.
	// An empty address slice and a non-nil error will be passed if the resolver returns
	// non-nil error to gRPC.
	//
	// Deprecated: if V2Balancer is implemented by the Balancer,
	// UpdateClientConnState will be called instead.
	HandleResolvedAddrs([]resolver.Address, error)
	// Close closes the balancer. The balancer is not required to call
	// ClientConn.RemoveSubConn for its existing SubConns.
	Close()
}
```
3. pick 接口
```go
// V2Picker is used by gRPC to pick a SubConn to send an RPC.
// Balancer is expected to generate a new picker from its snapshot every time its
// internal state has changed.
//
// The pickers used by gRPC can be updated by ClientConn.UpdateBalancerState().
type V2Picker interface {
	// Pick returns the connection to use for this RPC and related information.
	//
	// Pick should not block.  If the balancer needs to do I/O or any blocking
	// or time-consuming work to service this call, it should return
	// ErrNoSubConnAvailable, and the Pick call will be repeated by gRPC when
	// the Picker is updated (using ClientConn.UpdateState).
	//
	// If an error is returned:
	//
	// - If the error is ErrNoSubConnAvailable, gRPC will block until a new
	//   Picker is provided by the balancer (using ClientConn.UpdateState).
	//
	// - If the error implements IsTransientFailure() bool, returning true,
	//   wait for ready RPCs will wait, but non-wait for ready RPCs will be
	//   terminated with this error's Error() string and status code
	//   Unavailable.
	//
	// - Any other errors terminate all RPCs with the code and message
	//   provided.  If the error is not a status error, it will be converted by
	//   gRPC to a status error with code Unknown.
	Pick(info PickInfo) (PickResult, error)
}
```

4. 最后通过Pick获取连接并执行 【 call -> pick -> sendMsg -> recvMsg】 grpc 实际性调用

