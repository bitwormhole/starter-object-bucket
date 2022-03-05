package buckets

import "io"

// Connection 表示与某个具体存储桶的连接
type Connection interface {
	io.Closer

	Check() error

	GetObject(name string) Object
}
