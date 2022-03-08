package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// OSSDriver ...
type OSSDriver struct {
	markup.Component `class:"buckets.Driver" initMethod:"Init"`

	connector aliyunConnector
}

func (inst *OSSDriver) _Impl() (buckets.DriverRegistry, buckets.Driver) {
	return inst, inst
}

// Init ...
func (inst *OSSDriver) Init() error {
	return nil
}

// ListDrivers ...
func (inst *OSSDriver) ListDrivers() []*buckets.DriverRegistration {

	vlog.Info("aliyun.oss.version=", oss.Version)

	dr := &buckets.DriverRegistration{}
	dr.Name = "aliyun"
	dr.Driver = inst
	return []*buckets.DriverRegistration{dr}
}

// GetBucket ...
func (inst *OSSDriver) GetBucket(tag, id string, p collection.Properties) (*buckets.Bucket, error) {
	ldr := core.BucketLoader{}
	ldr.WantBucketExt = []string{pBucketEndpoint, pBucketName}
	ldr.WantCredentialExt = []string{pAccessKeyID, pAccessKeySecret}
	return ldr.Load(tag, id, p)
}

// GetConnector ...
func (inst *OSSDriver) GetConnector() buckets.Connector {
	return &inst.connector
}
