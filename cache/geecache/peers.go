package geecache

import pb "geecache/geecachepb"

// PeerPicker 节点选择
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从对应 group 查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
