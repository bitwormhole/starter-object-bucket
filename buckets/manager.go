package buckets

import "github.com/bitwormhole/starter/collection"

// Manager 用来统一管理各种驱动的管理器
// 【inject:"#buckets.Manager"】
type Manager interface {
	FindDriver(name string) (Driver, error)

	OpenBucket(b *Bucket) (Connection, error)

	GetBucket(tag, id string, p collection.Properties) (*Bucket, error)

	// // 根据 bucket(名称),profile(类型),zone 计算域名
	// ComputeDomainName(b *Bucket, profile Profile) (string, error)
}
