package geecache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "234",
	"Jack": "567",
	"Zee":  "11",
	"test": "test",
}

func Test_GeeCache(t *testing.T) {
	// 表示被回调的次数
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("getterFunc")
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key not found")
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal(err)
		}
		if _, err := gee.Get(k); err != nil {
			t.Fatal("cache key not found")
		}
	}
	if view, err := gee.Get("test"); err != nil || view.String() != "test" {
		t.Fatalf("cache key[%s] not found", "test")
	}
}
