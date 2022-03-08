package qq

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSDriver ...
type COSDriver struct {
	markup.Component `class:"buckets.Driver" initMethod:"Init"`

	connector cosConnector
}

func (inst *COSDriver) _Impl() (buckets.DriverRegistry, buckets.Driver) {
	return inst, inst
}

// Init ...
func (inst *COSDriver) Init() error {
	return nil
}

// ListDrivers ...
func (inst *COSDriver) ListDrivers() []*buckets.DriverRegistration {

	vlog.Info("qq.cos.version=", cos.Version)

	dr := &buckets.DriverRegistration{}
	dr.Name = "qq"
	dr.Driver = inst
	return []*buckets.DriverRegistration{dr}
}

// GetBucket ...
func (inst *COSDriver) GetBucket(tag, id string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := core.BucketLoader{}
	ldr.WantBucketExt = []string{pBucketEndpoint, pBucketName, pServiceURL, pBucketURL}
	ldr.WantCredentialExt = []string{pBucketAK, pBucketSK}
	return ldr.Load(tag, id, p)
}

// GetConnector ...
func (inst *COSDriver) GetConnector() buckets.Connector {
	return &inst.connector
}
