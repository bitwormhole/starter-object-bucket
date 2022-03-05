package aliyun

import (
	"errors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/bitwormhole/starter-object-bucket/buckets"
)

// OSS 的扩展参数
const (

	// for 桶

	BucketName     = "bucket"
	BucketEndpoint = "endpoint"

	// for 凭证

	BucketAccessKeyID     = "access-key-id"
	BucketAccessKeySecret = "access-key-secret"
)

type ossBucketConnection struct {
	bucketName string
	client     *oss.Client
	bucket     *oss.Bucket
}

func (inst *ossBucketConnection) _Impl() buckets.Connection {
	return inst
}

func (inst *ossBucketConnection) init(b *buckets.Bucket) error {

	ext := b.Ext
	endpoint := ext[BucketEndpoint]
	akeyID := ext[BucketAccessKeyID]
	akeySecret := ext[BucketAccessKeySecret]
	bucketName := ext[BucketName]

	client, err := oss.New(endpoint, akeyID, akeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	inst.bucketName = bucketName
	inst.client = client
	inst.bucket = bucket
	return nil
}

func (inst *ossBucketConnection) Close() error {
	// NOP
	return nil
}

func (inst *ossBucketConnection) Check() error {
	name := inst.bucketName
	ext, err := inst.client.IsBucketExist(name)
	if err != nil {
		return err
	}
	if !ext {
		return errors.New("the bucket is not exists, name:" + name)
	}
	return nil
}

func (inst *ossBucketConnection) GetObject(name string) buckets.Object {
	o := &ossBucketObject{
		parent: inst,
		name:   name,
	}
	return o
}

////////////////////////////////////////////////////////////////////////////////

type ossBucketObject struct {
	parent *ossBucketConnection
	name   string
}

func (inst *ossBucketObject) _Impl() buckets.Object {
	return inst
}

func (inst *ossBucketObject) Exists() bool {
	return false
}

func (inst *ossBucketObject) GetDownloadURL() string {
	return ""
}

func (inst *ossBucketObject) GetMeta() *buckets.ObjectMeta {
	return nil
}

func (inst *ossBucketObject) GetName() string {
	return ""
}

func (inst *ossBucketObject) UpdateMeta(meta *buckets.ObjectMeta) error {
	return errors.New("no impl")
}

func (inst *ossBucketObject) GetEntity() (buckets.ObjectEntity, error) {
	return nil, errors.New("no impl")
}

func (inst *ossBucketObject) PutEntity(entity buckets.ObjectEntity, meta *buckets.ObjectMeta) error {
	return errors.New("no impl")
}
