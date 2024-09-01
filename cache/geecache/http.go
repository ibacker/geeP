package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

// HTTPPool 通信节点池
type HTTPPool struct {
	// 节点地址
	self string
	// 节点通信前缀
	basePath string
}

// NewHTTPPool 实例化一个通信节点
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, args ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, args...))
}

// 在一个通信节点中进行缓存
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool.ServerHTTP: invalid path: " + r.URL.Path)
	}

	p.Log("HTTP %s %s", r.Method, r.URL.Path)

	// 拆分通信路径
	parts := strings.Split(r.URL.Path[len(p.basePath):], "/")
	if len(parts) != 2 {
		http.Error(w, "invalid path: "+r.URL.Path, http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	// 缓存 group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "invalid group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, "get group "+groupName+": "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
