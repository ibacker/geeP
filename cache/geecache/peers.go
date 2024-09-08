package geecache

// PeerPicker 节点选择
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从对应 group 查找缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
