package geecache

import (
	"fmt"
	pb "geecache/geecachepb"
	"geecache/singleflight"
	"log"
	"sync"
)

// Getter 定义接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 定义函数类型 -- 接口型函数
type GetterFunc func(key string) ([]byte, error)

// Get ： GetterFunc 实现 Getter 接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 缓存命名空间
type Group struct {
	name string
	// 缓存未命中时的回调函数
	getter    Getter
	mainCache cache
	// 节点
	peers PeerPicker
	// 避免 key 多次取值
	loader *singleflight.Group
}

var (
	// 新增 group 共享该读写锁
	mu sync.RWMutex
	// 维护多个缓存 group
	groups = make(map[string]*Group)
)

// NewGroup 新增 group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},

		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 获取缓存 group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 从缓存命名空间中获取对应值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache Get] key:", key, "value:", v)
		return v, nil
	}
	// 缓冲未命中，从远端获取
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 相同的 key 只会从远程调用1次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] failed o get from peer:", err)
			}
		}
		return g.getLocally(key)
	})

	if err != nil {
		return viewi.(ByteView), err
	}
	return
}

// 调用用户回调函数获取数据，并设置缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache 添加缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers
func (g *Group) RegisterPeers(peers PeerPicker) {
	//避免重复调用
	if g.peers != nil {
		panic("RegisterPeers called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
