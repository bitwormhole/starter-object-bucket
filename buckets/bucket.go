package buckets

// Bucket 表示对某个具体存储桶的配置
type Bucket struct {
	ID         string
	Driver     string
	URL        string // 这是下载的基本URL
	Credential string
	Ext        map[string]string
}
