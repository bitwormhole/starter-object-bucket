package core

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/collection"
)

// 基本的 bucket 参数
const (
	ParamBucketID         = "id"
	ParamBucketDriver     = "driver"
	ParamBucketCredential = "credential"
	ParamBucketName       = "name"
	ParamBucketDN         = "dn-bucket"
	ParamEndpointDN       = "dn-endpoint"

	// ParamBucketPublicURL  = "url"
	// ParamBucketZone       = "zone"

)

// BucketLoader 是 Bucket 的加载器
type BucketLoader struct {
	WantBucketExt     []string // 扩展的 bucket 参数
	WantCredentialExt []string // 扩展的 Credential 参数
	WantDriverExt     []string // 扩展的 driver 参数
}

func (inst *BucketLoader) addWant(tag, id, name string, wants map[string]string) {
	key := tag + "." + id + "." + name
	wants[key] = name
}

func (inst *BucketLoader) addWants(tag, id string, namelist []string, wants map[string]string) {
	if namelist == nil {
		return
	}
	for _, name := range namelist {
		inst.addWant(tag, id, name, wants)
	}
}

func (inst *BucketLoader) listBaseParams() []string {
	return []string{
		ParamBucketCredential,
		ParamBucketDriver,
		ParamBucketName,
		ParamBucketDN,
		ParamEndpointDN,

		// ParamBucketPublicURL,
		// ParamBucketZone,
		// 以下两个参数由driver提供
		// ParamBucketDomainNameTemplate,
		// ParamEndpointDomainNameTemplate,
	}
}

// Load ...
func (inst *BucketLoader) Load(tag, id string, p collection.Properties) (*buckets.Bucket, error) {

	getter := p.Getter()
	wants := make(map[string]string) // map[fullkey] shortkey
	ext := make(map[string]string)
	credentialID := getter.GetString(tag+"."+id+".credential", "")
	driverID := getter.GetString(tag+"."+id+".driver", "")

	inst.addWants(tag, id, inst.listBaseParams(), wants)
	inst.addWants(tag, id, inst.WantBucketExt, wants)
	inst.addWants("credential", credentialID, inst.WantCredentialExt, wants)
	inst.addWants("bucket-driver", driverID, inst.WantDriverExt, wants)

	for fullkey, shortkey := range wants {
		ext[shortkey] = getter.GetString(fullkey, "")
	}
	err := getter.Error()
	if err != nil {
		return nil, err
	}

	b := &buckets.Bucket{}
	b.ID = id
	b.Provider = ext[ParamBucketDriver]
	b.Name = ext[ParamBucketName]
	b.Credential = ext[ParamBucketCredential]
	b.BucketDN = ext[ParamBucketDN]
	b.EndpointDN = ext[ParamEndpointDN]
	b.Ext = ext
	return b, nil
}

////////////////////////////////////////////////////////////////////////////////

// LoadBucketParams 从 properties 加载 bucket 参数
func LoadBucketParams(tag, name string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := BucketLoader{}
	return ldr.Load(tag, name, p)
}
