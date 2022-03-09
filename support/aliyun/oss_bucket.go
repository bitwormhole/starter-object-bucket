package aliyun

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/util"
)

// OSS 的扩展参数
const (

	// for 桶

	pBucketName     = "bucket"
	pBucketEndpoint = "endpoint"
	pBucketBaseURL  = "url"

	// for 凭证

	pAccessKeyID     = "access-key-id"
	pAccessKeySecret = "access-key-secret"
)

// 对象大小界限
const (
	maxMiddleSize = 4 * 1024 * 1024 * 1024
	minMiddleSize = 8 * 1024 * 1024
)

////////////////////////////////////////////////////////////////////////////////

type ossBucketConnection struct {
	bucketName   string
	client       *oss.Client
	bucket       *oss.Bucket
	fetchBaseURL string
}

func (inst *ossBucketConnection) _Impl() buckets.Connection {
	return inst
}

func (inst *ossBucketConnection) init(b *buckets.Bucket) error {

	ext := b.Ext
	endpoint := ext[pBucketEndpoint]
	akeyID := ext[pAccessKeyID]
	akeySecret := ext[pAccessKeySecret]
	bucketName := ext[pBucketName]
	bucketURL := b.URL

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
	inst.fetchBaseURL = bucketURL
	return nil
}

func (inst *ossBucketConnection) Close() error {
	// NOP
	return nil
}

func (inst *ossBucketConnection) Check() error {

	// name := inst.bucketName
	// ext, err := inst.client.IsBucketExist(name)
	// if err != nil {
	// 	return err
	// }
	// if !ext {
	// 	return errors.New("the bucket is not exists, name:" + name)
	// }

	return nil
}

func (inst *ossBucketConnection) GetObject(name string) buckets.Object {
	o := &ossObject{
		parent: inst,
		name:   name,
	}
	return o
}

////////////////////////////////////////////////////////////////////////////////

type ossObject struct {
	parent *ossBucketConnection
	name   string
}

func (inst *ossObject) _Impl() buckets.Object {
	return inst
}

func (inst *ossObject) Exists() (bool, error) {
	return inst.parent.bucket.IsObjectExist(inst.name)
}

func (inst *ossObject) GetDownloadURL() string {
	p1 := inst.parent.fetchBaseURL
	p2 := inst.name
	return core.ComputeDownloadURL(p1, p2)
}

func (inst *ossObject) GetName() string {
	return inst.name
}

func (inst *ossObject) GetMeta() (*buckets.ObjectMeta, error) {
	bucket := inst.parent.bucket
	name := inst.name
	dst := map[string]string{}
	src, err := bucket.GetObjectDetailedMeta(name)
	if err != nil {
		return nil, err
	}
	for k := range src {
		k2 := strings.ToLower(k)
		dst[k2] = src.Get(k)
	}
	date, _ := time.Parse(time.RFC1123, dst["date"])
	size, _ := strconv.ParseInt(dst["content-length"], 10, 64)
	md5sum := util.Base64FromString(dst["content-md5"])
	meta := &buckets.ObjectMeta{}
	meta.More = dst
	meta.ContentType = dst["content-type"]
	meta.Hash = md5sum.HexString().String()
	meta.HashAlgorithm = "MD5"
	meta.Size = size
	meta.Date = date
	return meta, nil
}

func (inst *ossObject) UpdateMeta(meta *buckets.ObjectMeta) error {
	if meta == nil {
		return nil
	}
	src := meta.More
	if src == nil {
		return nil
	}
	options := []oss.Option{}
	for k, v := range src {
		item := oss.Meta(k, v)
		options = append(options, item)
	}
	bucket := inst.parent.bucket
	return bucket.SetObjectMeta(inst.name, options...)
}

func (inst *ossObject) GetEntity() (buckets.ObjectEntity, error) {
	return nil, errors.New("no impl")
}

// 根据文件大小计算 part_size
func (inst *ossObject) computePartSize(file fs.Path) int64 {
	const maxPartCount = 10000                 // 最多10000块
	const maxPartSize = 4 * 1024 * 1024 * 1024 // 每一块最大4GB
	total := file.Size()
	partSize := 8 * int64(1024*1024) // 完美的块大小大约为 8MB
	partCount := total / partSize    // 总块数
	for partCount > maxPartCount {
		partCount /= 2
		partSize *= 2 // 增大每一块的大小
	}
	if partSize > maxPartSize {
		partSize = maxPartSize
	}
	return partSize
}

func (inst *ossObject) PutFile(file fs.Path, meta *buckets.ObjectMeta) error {

	bucket := inst.parent.bucket
	src := file.Path()
	dst := inst.name
	partSize := inst.computePartSize(file) // 100*1024
	op1 := oss.Routines(3)                 // 并发数量
	op2 := oss.Checkpoint(true, "")        // 启用断点续传

	return bucket.UploadFile(dst, src, partSize, op1, op2)
}

func (inst *ossObject) PutEntity(entity buckets.ObjectEntity, meta *buckets.ObjectMeta) error {
	up := inst.getUploader(entity)
	return up.upload(meta, entity)
}

func (inst *ossObject) getUploader(entity buckets.ObjectEntity) uploader {
	size := entity.GetSize()
	if size < minMiddleSize {
		return &smallUploader{object: inst}
	} else if size < maxMiddleSize {
		return &middleUploader{object: inst}
	}
	return &largeUploader{object: inst}
}

////////////////////////////////////////////////////////////////////////////////

type uploader interface {
	upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error
}

////////////////////////////////////////////////////////////////////////////////

// 用文件缓冲，分多部上传
type largeUploader struct {
	object *ossObject
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
	object *ossObject
}

func (inst *middleUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {
	// 直接复用 smallUploader
	up := smallUploader{object: inst.object}
	return up.upload(meta, entity)
}

////////////////////////////////////////////////////////////////////////////////

// 用RAM缓冲，简单上传
type smallUploader struct {
	object *ossObject
}

func (inst *smallUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {

	bucket := inst.object.parent.bucket
	objectName := inst.object.name

	src, err := entity.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	err = bucket.PutObject(objectName, src)
	if err != nil {
		return err
	}

	// logger := vlog.Default()
	// if logger.IsDebugEnabled() {
	// 	logger.Debug("upload ", etag, " ... done")
	// }

	return nil
}

////////////////////////////////////////////////////////////////////////////////
