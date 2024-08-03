package geecache

// ByteView 只读数据结构，用于表示缓存值
// 实现了 Len 接口
type ByteView struct {
	// 存储真实的缓存值
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回复制值
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// 转换为字符串
func (v ByteView) String() string {
	return string(v.b)
}
