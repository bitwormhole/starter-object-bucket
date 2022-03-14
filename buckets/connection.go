package buckets

import (
	"io"
)

// // Profile 定义域名类型 (属于是flag类型，可以通过or操作组合)
// type Profile int

// // 定义各种域名类型
// const (
// 	ProfileMaskScope  Profile = 0xff00
// 	ProfileMaskAccess Profile = 0x00ff

// 	// for scope (bucket | endpoint)
// 	ProfileEndpoint Profile = 0x100
// 	ProfileBucket   Profile = 0x200

// 	// for access
// 	ProfileAcc      Profile = 0x01
// 	ProfileCustomer Profile = 0x02
// 	ProfileInternal Profile = 0x04
// 	ProfilePublic   Profile = 0x08
// 	ProfileVPC      Profile = 0x10
// )

// func (p Profile) String() string {
// 	n := int(p)
// 	return strconv.Itoa(n)
// }

////////////////////////////////////////////////////////////////////////////////

// Connection 表示与某个具体存储桶的连接
type Connection interface {
	io.Closer

	Check() error

	GetObject(name string) Object

	GetBucketName() string

	// GetDomainName(profile Profile) (string, error)
}
