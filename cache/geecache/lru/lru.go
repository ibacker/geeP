package lru

import "container/list"

// Cache LRU 缓存，非线程安全
type Cache struct {
	// 最大使用内存
	maxBytes int64
	// 已使用内存
	nbytes int64
	// 双向链表
	ll    *list.List
	cache map[string]*list.Element
	// 缓存驱逐
	onEvicted func(key string, value Value)
}

// 双向链表存储的对象
type entry struct {
	key   string
	value Value
}

// Value 被缓存对象，需实现该接口用于计算对象大小
type Value interface {
	Len() int
}

// New 初始化
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get 从缓存中获取元素，将链表中对应节点移动到头部
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 从缓存双向链表中获取元素 c.cache
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除元素-缓存驱逐
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从缓存字典中删除该节点的映射关系
		delete(c.cache, kv.key)
		c.nbytes -= int64(kv.value.Len()) + int64(len(kv.key))
		// 缓存驱逐事件
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Add 向缓存中添加元素，如果 key 已使用，替换原有缓存对象
func (c *Cache) Add(key string, value Value) {
	// 缓存中已存在
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		// (*entry) 类型断言 将元素转换为 entry 结构体
		kv := ele.Value.(*entry)
		// 使用新对象替换 旧对象，对应内存大小计算
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 缓存中不存在
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 缓存驱逐
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
