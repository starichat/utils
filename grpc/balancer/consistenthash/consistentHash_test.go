package consistenthash

import (
	"github.com/spaolacci/murmur3"
	"strconv"
	"testing"
)

func TestConsistentHash_Get(t *testing.T) {
	consistenhashIns := NewConsistentHash(6400, func(data []byte) uint64 {
		return murmur3.Sum64(data)
	})
	//添加物理节点
	consistenhashIns.Add("node-1")
	consistenhashIns.Add("node-2")
	consistenhashIns.Add("node-3")
	consistenhashIns.Add("node-4")
	consistenhashIns.Add("node-5")

	//测试请求
	for i := 0; i < 100; i++ {
		key := "request" + strconv.Itoa(i)
		node, ok := consistenhashIns.Get(key)
		t.Logf("当前请求：%s,哈希映射的物理节点:%s,%v", key, node, ok)
	}
	//移除node-1节点
	consistenhashIns.Remove("node-3")
	//再次测试请求
	t.Log("------second--------")
	for i := 0; i < 100; i++ {
		key := "request" + strconv.Itoa(i)
		node, ok := consistenhashIns.Get(key)
		t.Logf("当前请求：%s,哈希映射的物理节点:%s,%v", key, node, ok)
	}

}
