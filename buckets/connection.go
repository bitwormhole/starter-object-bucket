package buckets

import "io"

type DomainType string

// 定义各种域名类型
const (
	BucketAcc      DomainType = "bucket-acc-dn"
	BucketCustomer DomainType = "bucket-customer-dn"
	BucketInternal DomainType = "bucket-internal-dn"
	BucketPublic   DomainType = "bucket-public-dn"
	BucketVPC      DomainType = "bucket-vpc-dn"

	EndpointAcc      DomainType = "endpoint-acc-dn"
	EndpointCustomer DomainType = "endpoint-customer-dn"
	EndpointInternal DomainType = "endpoint-internal-dn"
	EndpointPublic   DomainType = "endpoint-public-dn"
	EndpointVPC      DomainType = "endpoint-vpc-dn"
)

// Connection 表示与某个具体存储桶的连接
type Connection interface {
	io.Closer

	Check() error

	GetObject(name string) Object

	GetBucketName() string

	GetDomainName(dntype DomainType) (string, error)
}
