package qq

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// bucket 参数
const (
	pBucketURL  = "bucket-url"
	pServiceURL = "service-url"

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

// 元数据键名
const (
	metaContentType   = "content-type"     // like 'application/octet-stream'
	metaContentLength = "content-length"   // like '16807'
	metaETag          = "etag"             // like '9a4802d5c99dafe1c04da0a8e7e166bf'
	metaLastModified  = "last-modified"    // like 'Wed, 28 Oct 2014 20:30:00 GMT'
	metaXCOSRequestID = "x-cos-request-id" // like 'NTg3NzQ3ZmVfYmRjMzVfMzE5N182NzczMQ=='
)

////////////////////////////////////////////////////////////////////////////////

type cosBucket struct {
	context      context.Context
	client       *cos.Client
	bucketName   string // the bucket name
	fetchBaseURL string
}

func (inst *cosBucket) _Impl() buckets.Connection {
	return inst
}

func (inst *cosBucket) getContext() context.Context {
	ctx := inst.context
	if ctx == nil {
		ctx = context.Background()
		inst.context = ctx
	}
	return ctx
}

func (inst *cosBucket) init(b *buckets.Bucket) error {

	ext := b.Ext
	ak := ext[pBucketAK]
	sk := ext[pBucketSK]
	bucketURL := ext[pBucketURL]
	serviceURL := ext[pServiceURL]
	bName := ext[pBucketName]

	bu, _ := url.Parse(bucketURL)
	su, _ := url.Parse(serviceURL)
	baseURL := &cos.BaseURL{BucketURL: bu, ServiceURL: su}

	transport := &cos.AuthorizationTransport{
		SecretID:  ak,
		SecretKey: sk,
	}
	httpclient := &http.Client{
		Transport: transport,
	}
	client := cos.NewClient(baseURL, httpclient)

	inst.bucketName = bName
	inst.client = client
	inst.fetchBaseURL = b.URL

	return inst.Check()
}

func (inst *cosBucket) Close() error {
	return nil // NOP
}

func (inst *cosBucket) Check() error {
	ctx := inst.getContext()
	ok, err := inst.client.Bucket.IsExist(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("the bucket is not exists, name=" + inst.bucketName)
	}
	return nil
}

func (inst *cosBucket) GetObject(name string) buckets.Object {
	o := &cosObject{
		parent: inst,
		name:   name,
	}
	return o
}

////////////////////////////////////////////////////////////////////////////////

type cosObject struct {
	context context.Context
	parent  *cosBucket
	name    string
}

func (inst *cosObject) _Impl() buckets.Object {
	return inst
}

func (inst *cosObject) getContext() context.Context {
	ctx := inst.context
	if ctx == nil {
		ctx = inst.parent.getContext()
		inst.context = ctx
	}
	return ctx
}

func (inst *cosObject) Exists() (bool, error) {
	ctx := inst.getContext()
	client := inst.parent.client
	name := inst.name
	return client.Object.IsExist(ctx, name)
}

func (inst *cosObject) GetDownloadURL() string {
	p1 := inst.parent.fetchBaseURL
	p2 := inst.name
	return core.ComputeDownloadURL(p1, p2)
}

func (inst *cosObject) GetMeta() (*buckets.ObjectMeta, error) {
	ctx := inst.getContext()
	client := inst.parent.client
	name := inst.name
	resp, err := client.Object.Head(ctx, name, nil)
	if err != nil {
		return nil, err
	}
	src := resp.Header
	dst := &buckets.ObjectMeta{}
	kvs := make(map[string]string)
	for key := range src {
		key2 := strings.ToLower(key)
		value := src.Get(key)
		kvs[key2] = value
	}
	date, _ := inst.parseDate(kvs[metaLastModified])
	size, _ := inst.parseSize(kvs[metaContentLength])
	etag := strings.ReplaceAll(kvs[metaETag], "\"", "")
	dst.Date = date
	dst.Size = size
	dst.Hash = strings.TrimSpace(etag)
	dst.HashAlgorithm = "MD5"
	dst.ContentType = kvs[metaContentType]
	dst.More = kvs
	return dst, nil
}

func (inst *cosObject) parseSize(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func (inst *cosObject) parseDate(s string) (time.Time, error) {
	const fmt = time.RFC1123
	return time.Parse(fmt, s)
}

func (inst *cosObject) GetName() string {
	return inst.name
}

func (inst *cosObject) GetEntity() (buckets.ObjectEntity, error) {
	return nil, errors.New("no impl")
}

func (inst *cosObject) UpdateMeta(meta *buckets.ObjectMeta) error {
	return errors.New("unsupported")
}

func (inst *cosObject) PutEntity(entity buckets.ObjectEntity, meta *buckets.ObjectMeta) error {
	up := inst.getUploader(entity)
	return up.upload(meta, entity)
}

func (inst *cosObject) PutFile(file fs.Path, meta *buckets.ObjectMeta) error {

	client := inst.parent.client
	objectName := inst.name
	filepath := file.Path()
	ctx := inst.getContext()

	// qq.cos 的这个方法，会根据用户文件的长度，自动切分数据，
	// 不需要区分简单上传和分片上传
	result, resp, err := client.Object.Upload(ctx, objectName, filepath, nil)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		if result != nil {
			logger.Debug("upload ", result.ETag)
		}
		if resp != nil {
			logger.Debug("upload ", resp.StatusCode)
		}
	}

	return nil
}

func (inst *cosObject) getUploader(entity buckets.ObjectEntity) uploader {
	size := entity.GetSize()
	const maxMiddleSize = 4 * 1024 * 1024 * 1024
	const minMiddleSize = 8 * 1024 * 1024
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
	object *cosObject
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
	object *cosObject
}

func (inst *middleUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {
	// 直接复用 smallUploader
	up := smallUploader{object: inst.object}
	return up.upload(meta, entity)
}

////////////////////////////////////////////////////////////////////////////////

// 用RAM缓冲，简单上传
type smallUploader struct {
	object *cosObject
}

func (inst *smallUploader) upload(meta *buckets.ObjectMeta, entity buckets.ObjectEntity) error {

	client := inst.object.parent.client
	objectName := inst.object.name
	ctx := inst.object.getContext()

	src, err := entity.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	resp, err := client.Object.Put(ctx, objectName, src, nil)
	if err != nil {
		return err
	}

	logger := vlog.Default()
	if logger.IsDebugEnabled() {
		if resp != nil {
			logger.Debug("upload ", resp.StatusCode)
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
