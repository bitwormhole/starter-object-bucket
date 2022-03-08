package baidu

import (
	"errors"
	"io/ioutil"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/util"
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

type bosBucket struct {
	client     *bos.Client
	bucketName string // the bucket name
}

func (inst *bosBucket) _Impl() buckets.Connection {
	return inst
}

func (inst *bosBucket) init(b *buckets.Bucket) error {

	ext := b.Ext
	ak := ext[pBucketAK]
	sk := ext[pBucketSK]
	endpoint := ext[pBucketEndpoint]
	bName := ext[pBucketName]

	clientConfig := bos.BosClientConfiguration{
		Ak:               ak,
		Sk:               sk,
		Endpoint:         endpoint,
		RedirectDisabled: false,
	}

	bosClient, err := bos.NewClientWithConfig(&clientConfig)
	if err != nil {
		return err
	}

	inst.bucketName = bName
	inst.client = bosClient
	return nil
}

func (inst *bosBucket) Close() error {
	return nil
}

func (inst *bosBucket) Check() error {
	ok, err := inst.client.DoesBucketExist(inst.bucketName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("the bucket is not exists, name=" + inst.bucketName)
	}
	return nil
}

func (inst *bosBucket) GetObject(name string) buckets.Object {
	o := &bosObject{
		parent: inst,
		name:   name,
	}
	return o
}

////////////////////////////////////////////////////////////////////////////////

type bosObject struct {
	parent *bosBucket
	name   string
}

func (inst *bosObject) _Impl() buckets.Object {
	return inst
}

func (inst *bosObject) Exists() (bool, error) {
	return false, errors.New("no impl")
}

func (inst *bosObject) GetDownloadURL() string {
	return ""
}

func (inst *bosObject) GetMeta() (*buckets.ObjectMeta, error) {
	return nil, errors.New("no impl")
}

func (inst *bosObject) GetName() string {
	return inst.name
}

func (inst *bosObject) GetEntity() (buckets.ObjectEntity, error) {
	return nil, errors.New("no impl")
}

func (inst *bosObject) UpdateMeta(meta *buckets.ObjectMeta) error {
	return errors.New("no impl")
}

func (inst *bosObject) PutFile(file fs.Path, meta *buckets.ObjectMeta) error {

	client := inst.parent.client
	bucket := inst.parent.bucketName
	obj := inst.name
	path := file.Path()

	res, err := client.ParallelUpload(bucket, obj, path, "", nil)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if res != nil && logger.IsDebugEnabled() {
		logger.Debug("upload ", res.ETag, " ... done")
	}

	return nil
}

func (inst *bosObject) PutEntity(entity buckets.ObjectEntity, meta *buckets.ObjectMeta) error {
	up := inst.getUploader(entity)
	return up.upload(meta, entity)
}

func (inst *bosObject) getUploader(entity buckets.ObjectEntity) uploader {
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
	object *bosObject
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
	object *bosObject
}

func (inst *middleUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {

	client := inst.object.parent.client
	bucketName := inst.object.parent.bucketName
	objectName := inst.object.name

	tmp := core.GetTempFileManager().NewTempFile()
	defer tmp.Close()
	file := tmp.GetPath()
	path := file.Path()

	err := inst.prepareTempFile(file, entity)
	if err != nil {
		return err
	}

	etag, err := client.PutObjectFromFile(bucketName, objectName, path, nil)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		logger.Debug("upload ", etag, " ... done")
	}

	return nil
}

func (inst *middleUploader) prepareTempFile(file fs.Path, entity buckets.ObjectEntity) error {

	opt := file.FileSystem().DefaultWriteOptions()
	opt.Create = true

	dst, err := file.GetIO().OpenWriter(opt, false)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := entity.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buffer := make([]byte, 64*1024)
	return util.PumpStream(src, dst, buffer)
}

////////////////////////////////////////////////////////////////////////////////

// 用RAM缓冲，简单上传
type smallUploader struct {
	object *bosObject
}

func (inst *smallUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {

	client := inst.object.parent.client
	bucketName := inst.object.parent.bucketName
	objectName := inst.object.name

	body, err := inst.makeBody(entity)
	if err != nil {
		return err
	}

	etag, err := client.PutObjectFromBytes(bucketName, objectName, body, nil)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		logger.Debug("upload ", etag, " ... done")
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
