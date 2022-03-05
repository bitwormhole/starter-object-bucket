package buckets

// Connector 表示一个连接器
type Connector interface {
	Open(b *Bucket) (Connection, error)
}
