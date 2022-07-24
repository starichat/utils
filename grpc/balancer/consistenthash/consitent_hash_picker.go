package consistenthash

import (
	"errors"
	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"log"
)

/**
pick 接口
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
*/

type Picker struct {
	consistentHash *ConsistentHash
	subConns       map[string]balancer.SubConn //string -> addr
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//从info中生成hashkey
	node, b := p.consistentHash.Get(GenerateKey(info))
	if !b {
		return balancer.PickResult{}, errors.New("no node")
	}
	if v, ok := p.subConns[string(node)]; ok {
		return balancer.PickResult{
			SubConn: v,
			Done:    nil,
		}, nil
	}
	return balancer.PickResult{}, errors.New("not found")
}

func GenerateKey(info balancer.PickInfo) string {
	key, ok := info.Ctx.Value("request_id").(string)
	if !ok {
		return uuid.New().String()
	}
	return key
}

func NewConsistentHashPicker(subConns map[string]balancer.SubConn) *Picker {
	addrs := make([]string, 0)
	for addr := range subConns {
		addrs = append(addrs, addr)
	}
	log.Printf("consistent hash picker built with addresses %v\n", addrs)
	hash := func(data []byte) uint64 {
		return murmur3.Sum64(data)
	}
	return &Picker{
		subConns:       subConns,
		consistentHash: NewConsistentHashWithAddrs(addrs, 65535, hash),
	}
}

type consistentHashPickerBuilder struct{}

func (b *consistentHashPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("consistentHashPicker: newPicker called with buildInfo: %v", buildInfo)
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	subConns := make(map[string]balancer.SubConn)
	for sc, conInfo := range buildInfo.ReadySCs {
		subConns[conInfo.Address.Addr] = sc
	}

	return NewConsistentHashPicker(subConns)
}
