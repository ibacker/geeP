package lru

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", String("1234"))
	if _, ok := lru.Get("key1"); !ok || lru.Len() != 1 {
		t.Fatalf("cache add key1 failed")
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("remove old key1 failed")
	}
}

func TestCache_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(0), callback)
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("1234"))
	lru.Add("key3", String("1234"))

	lru.RemoveOldest()
	lru.RemoveOldest()

	expected := []string{"key1", "key2"}
	if !reflect.DeepEqual(expected, keys) {
		t.Fatalf("cache onEvicted failed %s", keys)
	}
}

// sync.Mutex 的使用

var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
}

func Test_main(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce(i)
	}
	time.Sleep(time.Second)
}
