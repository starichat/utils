package consistenthash

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
