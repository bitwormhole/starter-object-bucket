package huawei

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

// OBSDriver ...
type OBSDriver struct {
	markup.Component `class:"buckets.Driver" initMethod:"Init"`

	BM buckets.Manager `inject:"#buckets.Manager"`

	connector obsConnector
}

func (inst *OBSDriver) _Impl() (buckets.DriverRegistry, buckets.Driver) {
	return inst, inst
}

// Init ...
func (inst *OBSDriver) Init() error {
	return nil
}

// ListDrivers ...
func (inst *OBSDriver) ListDrivers() []*buckets.DriverRegistration {

	vlog.Info("huawei.obs.version=", obs.USER_AGENT)

	dr := &buckets.DriverRegistration{}
	dr.Name = "huawei"
	dr.Driver = inst
	return []*buckets.DriverRegistration{dr}
}

// GetConnector ...
func (inst *OBSDriver) GetConnector() buckets.Connector {
	return &inst.connector
}

// GetBucket ...
func (inst *OBSDriver) GetBucket(tag, id string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := core.BucketLoader{}
	ldr.WantBucketExt = []string{pBucketEndpoint, pBucketName}
	ldr.WantCredentialExt = []string{pBucketAK, pBucketSK}
	return ldr.Load(tag, id, p)
}
