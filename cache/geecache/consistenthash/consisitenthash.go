package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 哈希算法 将 bytes 转换为 uint32
type Hash func(data []byte) uint32

type Map struct {
	// 哈希算法
	hash Hash
	// 虚拟节点
	replicas int
	// 节点列表
	keys []int
	// 缓存对象
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	// 默认哈希算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 节点的 get 方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 二分查找最近的节点 顺时针查找
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 取模
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// RemoveNode 删除节点
func (m *Map) RemoveNode(key string) {
	// 虚拟节点
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		// 删除虚拟节点
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		// 删除缓存
		delete(m.hashMap, hash)
	}

}
