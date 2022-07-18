package consistenthash

import (
	"sort"
	"strconv"
	"sync"
)

/**
一致性hash负载均衡
*/

type Func func(data []byte) uint64

//Node 物理节点代表的元素
type Node string

type ConsistentHash struct {
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

type Base interface {
	Add(node Node)
	Remove(node Node)
	Get(v string) (Node, bool)
}

func NewConsistentHash(size int, hash Func) *ConsistentHash {

	return &ConsistentHash{
		RWMutex:         sync.RWMutex{},
		hashFunc:        hash,
		ring:            make(map[uint64]Node),
		nodes:           make(map[Node]struct{}),
		NumVirtualNodes: size,
	}
}

//Add ...
func (h *ConsistentHash) Add(node Node) {
	h.Lock()
	defer h.Unlock()
	// 为每台服务器生成数量为 replicateCount-1 个虚拟节点
	// 并将其与服务器的实际节点一同添加到哈希环中
	for i := 0; i < h.NumVirtualNodes; i++ {
		// 获取节点的哈希值，其中节点的字符串为 node+index
		hashkey := h.hashFunc([]byte(Key(node, i)))
		h.ring[hashkey] = node
		// 将节点的哈希值添加到哈希环中
		h.keys = append(h.keys, hashkey)
		h.nodes[node] = struct{}{}
	}

	sort.Sort(h.keys)
}

//Remove ...
func (h *ConsistentHash) Remove(node Node) {
	h.Lock()
	defer h.Unlock()
	// 移除时需要将服务器的实际节点和虚拟节点一同移除
	for i := 0; i < h.NumVirtualNodes; i++ {
		// 计算节点的哈希值
		hashkey := h.hashFunc([]byte(Key(node, i)))

		//找到虚拟节点索引位置
		index := sort.Search(len(h.keys), func(i int) bool {
			return h.keys[i] >= hashkey
		})

		if index < len(h.keys) && h.keys[index] == hashkey {
			//移除index位置的元素
			h.keys = append(h.keys[:index], h.keys[index+1:]...)
		}
		// 移除虚拟节点到物理节点的映射关系
		delete(h.ring, hashkey)
	}
	//移除真实节点
	delete(h.nodes, node)
}

//Get ...
func (h *ConsistentHash) Get(key string) (Node, bool) {
	h.Lock()
	defer h.Unlock()
	if len(h.ring) == 0 || len(h.nodes) == 0 {
		return "", false
	}

	// 获取客户端地址的哈希值
	hashKey := h.hashFunc([]byte(key))
	//确定索引位置
	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hashKey
	}) % len(h.keys)

	//通过上述二分查找算法找到对应的物理节点
	node := h.ring[h.keys[index]]
	_, ok := h.nodes[node]
	return node, ok
}
