package buckets

// Bucket 表示对某个具体存储桶的配置
type Bucket struct {
	ID       string
	Name     string // 桶名称
	Provider string // 驱动名称

	EndpointDN string
	BucketDN   string
	Credential string

	Driver Driver
	Ext    map[string]string
}

////////////////////////////////////////////////////////////////////////////////

// // DomainNameBuilder 域名创建器
// type DomainNameBuilder struct {
// 	BucketName string
// 	Zone       string
// 	Profile    string
// 	Template   string
// }

// // DomainName 创建域名
// func (inst *DomainNameBuilder) DomainName() string {

// 	dn := inst.Template
// 	name := inst.BucketName
// 	zone := inst.Zone
// 	profile := inst.Profile

// 	dn = strings.ReplaceAll(dn, "{{zone}}", zone)
// 	dn = strings.ReplaceAll(dn, "{{name}}", name)
// 	dn = strings.ReplaceAll(dn, "{{profile}}", profile)

// 	return dn
// }
