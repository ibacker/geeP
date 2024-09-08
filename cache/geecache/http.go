package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_geecache/"
const defaultReplicas = 50

// HTTPPool 通信节点池
type HTTPPool struct {
	// 节点地址
	self string
	// 节点通信前缀
	basePath string
	mu       sync.Mutex
	// 根据具体的 key 选择节点
	peers *consistenthash.Map
	// 映射远程节点与对应的 httpGetter
	httpGetters map[string]*httpGetter
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

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("PickPeer: peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

// 客户端功能

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
