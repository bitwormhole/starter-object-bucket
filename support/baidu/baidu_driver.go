package baidu

import (
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// BOSDriver ...
type BOSDriver struct {
	markup.Component `class:"buckets.Driver" initMethod:"Init"`

	BM buckets.Manager `inject:"#buckets.Manager"`

	connector bosConnector
}

func (inst *BOSDriver) _Impl() (buckets.DriverRegistry, buckets.Driver) {
	return inst, inst
}

// Init ...
func (inst *BOSDriver) Init() error {
	return nil
}

// ListDrivers ...
func (inst *BOSDriver) ListDrivers() []*buckets.DriverRegistration {

	v := bos.DEFAULT_SERVICE_DOMAIN
	vlog.Info("baidu.bos.version=...", v)

	dr := &buckets.DriverRegistration{}
	dr.Name = "baidu"
	dr.Driver = inst
	return []*buckets.DriverRegistration{dr}
}

// GetBucket ...
func (inst *BOSDriver) GetBucket(tag, name string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := core.BucketLoader{}
	ldr.WantBucketExt = []string{pBucketEndpoint, pBucketName}
	ldr.WantCredentialExt = []string{pBucketAK, pBucketSK}
	return ldr.Load(tag, name, p)
}

// GetConnector ...
func (inst *BOSDriver) GetConnector() buckets.Connector {
	return &inst.connector
}
