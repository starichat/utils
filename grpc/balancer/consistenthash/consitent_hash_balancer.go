package consistenthash

import (
	"google.golang.org/grpc/balancer"
	"sync"
)

/**
builder接口
// Builder creates a balancer.
type Builder interface {
	// Build creates a new balancer with the ClientConn.
	Build(cc ClientConn, opts BuildOptions) Balancer
	// Name returns the name of balancers built by this builder.
	// It will be used to pick balancers (for example in service config).
	Name() string
}

*/

type BalancerBuilder struct {
	sync.RWMutex
	hashFunc        Func
	keys            slots                       //虚拟节点列表，排序好的虚拟节点，便于通过二分算法快速定位到最近的物理节点
	ring            map[uint64]balancer.SubConn //虚拟节点到物理节点的映射
	nodes           map[Node]struct{}           //物理节点映射，判断当前物理节点是否存在
	NumVirtualNodes int                         // 为每台机器在hash圆环上创建多少个虚拟Node
}

//实现Builder接口

//func newBalancerBuilder() BalancerBuilder {
//
//}

/**
Balancer接口

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
*/
