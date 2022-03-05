package baidu

import (
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

type BOSDriver struct {
	markup.Component `class:"buckets.Driver" initMethod:"Init"`

	connector bosConnector
}

func (inst *BOSDriver) _Impl() (buckets.DriverRegistry, buckets.Driver) {
	return inst, inst
}

func (inst *BOSDriver) Init() error {
	return nil
}

func (inst *BOSDriver) ListDrivers() []*buckets.DriverRegistration {

	v := bos.DEFAULT_SERVICE_DOMAIN
	vlog.Info("baidu.bos.version=...", v)

	dr := &buckets.DriverRegistration{}
	dr.Name = "baidu"
	dr.Driver = inst
	return []*buckets.DriverRegistration{dr}
}

func (inst *BOSDriver) GetBucket(tag, name string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := core.BucketLoader{}
	ldr.WantBucketExt = []string{BucketEndpoint, BucketName}
	ldr.WantCredentialExt = []string{BucketAK, BucketSK}
	return ldr.Load(tag, name, p)
}

func (inst *BOSDriver) GetConnector() buckets.Connector {
	return &inst.connector
}
