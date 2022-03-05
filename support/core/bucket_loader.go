package core

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/collection"
)

// 基本的 bucket 参数
const (
	BaseBucketParamID         = "id"
	BaseBucketParamURL        = "url"
	BaseBucketParamDriver     = "driver"
	BaseBucketParamCredential = "credential"
)

// BucketLoader 是 Bucket 的加载器
type BucketLoader struct {
	WantBucketExt     []string // 扩展的 bucket 参数
	WantCredentialExt []string // 扩展的 Credential 参数
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
	return []string{BaseBucketParamCredential, BaseBucketParamDriver, BaseBucketParamID, BaseBucketParamURL}
}

// Load ...
func (inst *BucketLoader) Load(tag, id string, p collection.Properties) (*buckets.Bucket, error) {

	getter := p.Getter()
	wants := make(map[string]string) // map[fullkey] shortkey
	ext := make(map[string]string)
	credentialID := getter.GetString(tag+"."+id+".credential", "")

	inst.addWants(tag, id, inst.listBaseParams(), wants)
	inst.addWants(tag, id, inst.WantBucketExt, wants)
	inst.addWants("credential", credentialID, inst.WantCredentialExt, wants)

	for fullkey, shortkey := range wants {
		ext[shortkey] = getter.GetString(fullkey, "")
	}
	err := getter.Error()
	if err != nil {
		return nil, err
	}

	b := &buckets.Bucket{}
	b.Credential = ext[BaseBucketParamCredential]
	b.Driver = ext[BaseBucketParamDriver]
	b.ID = ext[BaseBucketParamID]
	b.URL = ext[BaseBucketParamURL]
	b.Ext = ext
	return b, nil
}

////////////////////////////////////////////////////////////////////////////////

// LoadBucketParams 从 properties 加载 bucket 参数
func LoadBucketParams(tag, name string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := BucketLoader{}
	return ldr.Load(tag, name, p)
}
