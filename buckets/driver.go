package buckets

import "github.com/bitwormhole/starter/collection"

// Driver 表示为某种存储方案提供的的驱动
type Driver interface {
	GetConnector() Connector

	// 从 properties 中加载 bucket 的参数
	GetBucket(tag, id string, p collection.Properties) (*Bucket, error)

	// StringifyProfile(p Profile) string
}

// DriverRegistration 用来表示某个驱动的注册信息
type DriverRegistration struct {
	Name   string
	Driver Driver
}

// DriverRegistry 用来枚举供应商提供的驱动
// 【inject:".buckets.Driver"】
type DriverRegistry interface {
	ListDrivers() []*DriverRegistration
}
