package consistenthash

import (
	"github.com/spaolacci/murmur3"
	"strconv"
	"testing"
)

func TestConsistentHash_Add(t *testing.T) {
	consistenhashIns := NewConsistentHash(64, func(data []byte) uint64 {
		return murmur3.Sum64(data)
	})
	consistenhashIns.Add("node-1")
	consistenhashIns.Add("node-2")
	consistenhashIns.Add("node-3")

	//测试请求
	for i := 0; i < 20; i++ {
		key := "request" + strconv.Itoa(i)
		node, ok := consistenhashIns.Get(key)
		t.Logf("当前请求：%s,哈希映射的物理节点:%s,%v", key, node, ok)
	}

}
