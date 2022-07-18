package balancer

import (
	"strconv"
	"sync"
)

type Func func(data []byte) uint64

//Node 物理节点代表的元素
type Node string

type ConsistentHashGrpc struct {
	sync.RWMutex
	hashFunc        Func
	keys            slots             //虚拟节点列表，排序好的虚拟节点，便于通过二分算法快速定位到最近的物理节点
	ring            map[uint64]Node   //虚拟节点到物理节点的映射
	nodes           map[Node]struct{} //物理节点映射，判断当前物理节点是否存在
	NumVirtualNodes int               // 为每台机器在hash圆环上创建多少个虚拟Node
}

//Key ...
func Key(node Node, index int) string {
	return string(node) + strconv.Itoa(index)
}

// 使用sort.Sort函数，传入的参数需要实现的接口
type slots []uint64

func (s slots) Len() int {
	return len(s)
}

func (s slots) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s slots) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
