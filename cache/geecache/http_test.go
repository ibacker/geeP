package geecache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestHTTPPool_ServeHTTP(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("mock db search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not found", key)
		}))

	addr := "127.0.0.1:9999"
	peers := NewHTTPPool(addr)
	log.Println("geecache runs on", addr)
	http.ListenAndServe(addr, peers)

}
