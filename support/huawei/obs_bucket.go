package huawei

import (
	"errors"
	"io/ioutil"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

// bucket 参数
const (
	pBucketEndpoint = "endpoint"
	pBucketName     = "bucket"
	pBucketAK       = "access-key-id"
	pBucketSK       = "access-key-secret"
)

// 对象大小界限
const (
	maxMiddleSize = 4 * 1024 * 1024 * 1024
	minMiddleSize = 8 * 1024 * 1024
)

////////////////////////////////////////////////////////////////////////////////

type obsBucket struct {
	client     *obs.ObsClient
	bucketName string // the bucket name
}

func (inst *obsBucket) _Impl() buckets.Connection {
	return inst
}

func (inst *obsBucket) init(b *buckets.Bucket) error {

	ext := b.Ext
	ak := ext[pBucketAK]
	sk := ext[pBucketSK]
	endpoint := ext[pBucketEndpoint]
	bName := ext[pBucketName]

	client, err := obs.New(ak, sk, endpoint)
	if err != nil {
		return err
	}

	inst.bucketName = bName
	inst.client = client
	return nil
}

func (inst *obsBucket) Close() error {
	client := inst.client
	inst.client = nil
	if client != nil {
		client.Close()
	}
	return nil
}

func (inst *obsBucket) Check() error {
	_, err := inst.client.HeadBucket(inst.bucketName)
	return err
}

func (inst *obsBucket) GetObject(name string) buckets.Object {
	o := &obsObject{
		parent: inst,
		name:   name,
	}
	return o
}

func (inst *obsBucket) GetBucketName() string {
	return inst.bucketName
}

// func (inst *obsBucket) GetDomainName(p buckets.Profile) (string, error) {
// 	return "", errors.New("no impl")
// }

////////////////////////////////////////////////////////////////////////////////

type obsObject struct {
	parent *obsBucket
	name   string
}

func (inst *obsObject) _Impl() buckets.Object {
	return inst
}

func (inst *obsObject) Exists() (bool, error) {
	return false, errors.New("no impl")
}

func (inst *obsObject) GetDownloadURL() string {
	return ""
}

func (inst *obsObject) GetMeta() (*buckets.ObjectMeta, error) {
	return nil, errors.New("no impl")
}

func (inst *obsObject) GetName() string {
	return inst.name
}

func (inst *obsObject) GetEntity() (buckets.ObjectEntity, error) {
	return nil, errors.New("no impl")
}

func (inst *obsObject) UpdateMeta(meta *buckets.ObjectMeta) error {
	return errors.New("no impl")
}

func (inst *obsObject) PutEntity(entity buckets.ObjectEntity, meta *buckets.ObjectMeta) error {
	up := inst.getUploader(entity)
	return up.upload(meta, entity)
}

// 根据文件大小计算 part_size
func (inst *obsObject) computePartSize(file fs.Path) int64 {
	const maxPartCount = 10000                 // 最多10000块
	const maxPartSize = 4 * 1024 * 1024 * 1024 // 每一块最大4GB
	total := file.Size()
	partSize := 16 * int64(1024*1024) // 完美的块大小大约为16MB
	partCount := total / partSize     // 总块数
	for partCount > maxPartCount {
		partCount /= 2
		partSize *= 2 // 增大每一块的大小
	}
	if partSize > maxPartSize {
		partSize = maxPartSize
	}
	return partSize
}

func (inst *obsObject) PutFile(file fs.Path, meta *buckets.ObjectMeta) error {

	client := inst.parent.client

	input := &obs.UploadFileInput{}
	input.Bucket = inst.parent.bucketName
	input.Key = inst.name
	input.UploadFile = file.Path()
	input.EnableCheckpoint = true               // 开启断点续传模式
	input.PartSize = inst.computePartSize(file) // 指定分段大小为9MB
	input.TaskNum = 5                           // 指定分段上传时的最大并发数

	output, err := client.UploadFile(input)
	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		if err == nil && output != nil {
			logger.Debug("upload ", output.ETag, " ... done")
		}
	}
	return err
}

func (inst *obsObject) getUploader(entity buckets.ObjectEntity) uploader {
	size := entity.GetSize()
	if size < minMiddleSize {
		return &smallUploader{object: inst}
	} else if size < maxMiddleSize {
		return &middleUploader{object: inst}
	}
	return &largeUploader{object: inst}
}

func (inst *obsObject) UploadByAPI(up *buckets.HTTPUploading) (*buckets.HTTPUploading, error) {
	return nil, errors.New("no impl")
}

////////////////////////////////////////////////////////////////////////////////

type uploader interface {
	upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error
}

////////////////////////////////////////////////////////////////////////////////

// 用文件缓冲，分多部上传
type largeUploader struct {
	object *obsObject
}

func (inst *largeUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {
	tmp, err := core.PrepareLargeTempFileForUploading(entity)
	if err != nil {
		return err
	}
	defer tmp.Close()
	file := tmp.GetPath()
	return inst.object.PutFile(file, meta)
}

////////////////////////////////////////////////////////////////////////////////

// 用文件缓冲，简单上传
type middleUploader struct {
	object *obsObject
}

func (inst *middleUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {
	// 直接复用 smallUploader
	up := smallUploader{object: inst.object}
	return up.upload(meta, entity)
}

////////////////////////////////////////////////////////////////////////////////

// 用RAM缓冲，简单上传
type smallUploader struct {
	object *obsObject
}

func (inst *smallUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {

	client := inst.object.parent.client
	bucketName := inst.object.parent.bucketName
	objectName := inst.object.name

	src, err := entity.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	input := &obs.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectName
	input.Body = src
	output, err := client.PutObject(input)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		logger.Debug("upload ", output.ETag, " ... done")
	}

	return nil
}

func (inst *smallUploader) makeBody(entity buckets.ObjectEntity) ([]byte, error) {
	src, err := entity.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return ioutil.ReadAll(src)
}

////////////////////////////////////////////////////////////////////////////////
