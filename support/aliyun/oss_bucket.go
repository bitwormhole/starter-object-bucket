package aliyun

import (
	"errors"
	"net/http"
	"net/url"
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

	// pBucketEndpoint = "ext-bucket-endpoint"
	// pBucketName     = "ext-bucket-name"

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
	bucketName      string
	client          *oss.Client
	bucket          *oss.Bucket
	bucketDN        string
	fetchBaseURL    string
	dnSet           core.BucketDNSet
	accessKeyID     string
	accessKeySecret string
}

func (inst *ossBucketConnection) _Impl() buckets.Connection {
	return inst
}

// func (inst *ossBucketConnection) makeEndpointDN(b *buckets.Bucket) string {
// 	builder := buckets.DomainNameBuilder{}
// 	builder.Template = b.EndpointDomainTemplate
// 	builder.BucketName = b.Name
// 	builder.Profile = inst.stringifyAccess(b.Profile)
// 	builder.Zone = b.Zone
// 	return builder.DomainName()
// }

// func (inst *ossBucketConnection) stringifyAccess(p buckets.Profile) string {
// 	access := p & buckets.ProfileMaskAccess
// 	switch access {
// 	case buckets.ProfileAcc:
// 		return "-todo"
// 	case buckets.ProfileCustomer:
// 		return "-todo"
// 	case buckets.ProfileInternal:
// 		return "-internal"
// 	case buckets.ProfilePublic:
// 		return ""
// 	case buckets.ProfileVPC:
// 		return "-todo"
// 	default:
// 		return "-internal"
// 	}
// }

func (inst *ossBucketConnection) init(b *buckets.Bucket) error {

	ext := b.Ext
	endpoint := b.EndpointDN
	akeyID := ext[pAccessKeyID]
	akeySecret := ext[pAccessKeySecret]
	bucketName := b.Name
	bucketURL := "https://" + b.BucketDN

	client, err := oss.New(endpoint, akeyID, akeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	inst.bucketName = bucketName
	inst.bucketDN = b.BucketDN
	inst.client = client
	inst.bucket = bucket
	inst.fetchBaseURL = bucketURL
	inst.accessKeyID = akeyID
	inst.accessKeySecret = akeySecret
	inst.dnSet.Init(ext)
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

func (inst *ossBucketConnection) GetBucketName() string {
	return inst.bucketName
}

// func (inst *ossBucketConnection) GetDomainName(p buckets.Profile) (string, error) {
// 	return inst.dnSet.GetDN(p)
// }

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

func (inst *ossObject) UploadByAPI(up1 *buckets.HTTPUploading) (*buckets.HTTPUploading, error) {

	// 准备返回值
	up2 := &buckets.HTTPUploading{}
	headers2 := make(map[string]string)
	if up1 == nil {
		up1 = up2
	} else {
		*up2 = *up1
	}
	headers1 := up1.RequestHeaders
	for k, v := range headers1 {
		k2 := strings.ToLower(k)
		headers2[k2] = v
	}

	// 检查各个字段
	method := http.MethodPut
	path := inst.name
	contentType := up2.ContentType
	md5sum := up2.ContentMD5.String()
	length := up2.ContentLength
	dn := inst.parent.bucketDN

	if length < 0 {
		length = 0
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 写入 headers2
	options := []oss.Option{}
	options = append(options, oss.ContentType(contentType))
	headers2["content-type"] = contentType

	if md5sum != "" {
		options = append(options, oss.ContentMD5(md5sum))
	}

	if length > 0 {
		options = append(options, oss.ContentLength(length))
	}

	// 生成 URL
	url1, err := inst.parent.bucket.SignURL(path, oss.HTTPPut, 60, options...)
	if err != nil {
		return nil, err
	}

	url2, err := url.Parse(url1)
	if err != nil {
		return nil, err
	}

	if dn != "" {
		url2.Host = dn
	}

	if up1.UseHTTPS {
		url2.Scheme = "https"
	}

	up2.Method = method
	up2.URL = url2.String()
	up2.RequestHeaders = headers2
	up2.ContentLength = length
	up2.ContentMD5 = up1.ContentMD5
	up2.ContentType = contentType
	return up2, nil
}

////////////////////////////////////////////////////////////////////////////////

// type authorizationBuilder struct {
// 	VERB                    string // aka. HTTP.Method
// 	ContentMD5              string
// 	ContentType             string
// 	Date                    string
// 	CanonicalizedOSSHeaders string
// 	CanonicalizedResource   string

// 	AccessKeyId     string
// 	AccessKeySecret string
// }

// func (inst *authorizationBuilder) Create() string {

// 	const nl = ""
// 	aKeySecret := []byte(inst.AccessKeySecret)
// 	aKeyID := inst.AccessKeyId

// 	// 组装 value
// 	builder := strings.Builder{}
// 	builder.WriteString(inst.VERB)
// 	builder.WriteString(nl)

// 	builder.WriteString(inst.ContentMD5)
// 	builder.WriteString(nl)

// 	builder.WriteString(inst.ContentType)
// 	builder.WriteString(nl)

// 	builder.WriteString(inst.Date)
// 	builder.WriteString(nl)

// 	builder.WriteString(inst.CanonicalizedOSSHeaders)
// 	builder.WriteString(inst.CanonicalizedResource)
// 	value := builder.String()

// 	// 计算 hmac-sha1
// 	mac := hmac.New(sha1.New, aKeySecret)
// 	mac.Write([]byte(value))
// 	sum := mac.Sum(nil)

// 	// 生成签名
// 	signature := util.Base64FromBytes(sum)
// 	return "OSS " + aKeyID + ":" + signature.String()
// }

// func (inst *authorizationBuilder) init(bucketName, objectName string) {
// 	inst.ContentType = "application/octet-stream"
// 	inst.initCanonicalizedResource(bucketName, objectName)
// 	inst.initDate()
// }

// func (inst *authorizationBuilder) initCanonicalizedResource(bucketName, objectName string) {
// 	inst.CanonicalizedResource = "/" + bucketName + "/" + objectName
// }

// func (inst *authorizationBuilder) initDate() {
// 	now := time.Now()
// 	zone := time.FixedZone("GMT", 0)
// 	inst.Date = now.In(zone).Format(time.RFC1123)
// }

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
