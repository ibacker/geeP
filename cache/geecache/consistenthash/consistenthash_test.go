package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		// 转换为 int
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 添加三个真实缓存节点
	hash.Add("6", "4", "2")
	// 映射的虚拟缓存节点有  2 4 6 12 14 16 22 24 26
	// 真实节点 映射到虚拟节点
	// 6：6 16 26
	// 4： 4 14 24
	// 2： 2 12 22

	//
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 添加真实缓存节点 8
	hash.Add("8")
	// 映射虚拟节点 8：8 18 28

	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
