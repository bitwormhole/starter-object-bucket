package demo

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/util"
	"github.com/bitwormhole/starter/vlog"
)

// Demo3 测试 UploadByAPI
type Demo3 struct {
	markup.Component `id:"demo3"`

	DemoBuckets        string              `inject:"${demo.buckets}"`
	CredentialFileName string              `inject:"${demo.credential.properties}"`
	Context            application.Context `inject:"context"`
	BM                 buckets.Manager     `inject:"#buckets.Manager"`
}

func (inst *Demo3) _Impl() application.LifeRegistry {
	return inst
}

// GetLifeRegistration ...
func (inst *Demo3) GetLifeRegistration() *application.LifeRegistration {
	lr := &application.LifeRegistration{}
	lr.OnInit = inst.Init
	lr.Looper = inst
	return lr
}

// Init ...
func (inst *Demo3) Init() error {
	return inst.loadKeys()
}

// Loop ...
func (inst *Demo3) Loop() error {

	bucketlist, err := inst.loadBuckets()
	if err != nil {
		return err
	}

	for _, bucket := range bucketlist {
		err = inst.testUploadByAPIWithBucket(bucket)
		if err != nil {
			return err
		}
	}

	return nil
}

func (inst *Demo3) testUploadByAPIWithBucket(bucket *buckets.Bucket) error {

	driver, err := inst.BM.FindDriver(bucket.Driver)
	if err != nil {
		return err
	}

	props := inst.Context.GetProperties()
	bucket, err = driver.GetBucket("bucket", bucket.ID, props)
	if err != nil {
		return err
	}

	conn, err := driver.GetConnector().Open(bucket)
	if err != nil {
		return err
	}
	defer conn.Close()

	o1 := conn.GetObject("test/demo3/object1.txt")

	up1 := &buckets.HTTPUploading{
		UseHTTPS: true,
		Domain:   buckets.BucketPublic,
	}
	up2, err := o1.UploadByAPI(up1)
	if err != nil {
		return err
	}

	vlog.Debug("uploading: ", up2)
	text := "abcd1234." + time.Now().Format(time.RFC1123Z)
	return inst.doUpload(up2, []byte(text))
}

func (inst *Demo3) md5sum(data []byte) string {
	sum := md5.Sum(data)
	return util.StringifyBytes(sum[:])
}

func (inst *Demo3) doUpload(up *buckets.HTTPUploading, data []byte) error {

	md5sum := inst.md5sum(data)
	vlog.Info("[upload md5:", md5sum, " method:", up.Method, " url:", up.URL, "]")

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(up.Method, up.URL, body)
	if err != nil {
		return err
	}

	headers := up.RequestHeaders
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	code := resp.StatusCode
	body2 := resp.Body

	if body2 != nil {
		defer body2.Close()
	}

	body2data, err := ioutil.ReadAll(body2)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		vlog.Warn("response.body:", string(body2data))

		msg := strings.Builder{}
		msg.WriteString("HTTP ")
		msg.WriteString(resp.Status)
		return errors.New(msg.String())
	}

	return nil
}

func (inst *Demo3) loadBuckets() ([]*buckets.Bucket, error) {
	dst := make([]*buckets.Bucket, 0)
	items := strings.Split(inst.DemoBuckets, ",")
	for _, item := range items {
		b, err := inst.loadBucketWithName(item)
		if err != nil {
			return nil, err
		}
		dst = append(dst, b)
	}
	return dst, nil
}

func (inst *Demo3) loadBucketWithName(name string) (*buckets.Bucket, error) {

	name = strings.TrimSpace(name)

	b := &buckets.Bucket{}
	p := "bucket." + name + "."
	getter := inst.Context.GetProperties().Getter()

	b.Driver = getter.GetString(p+"driver", "")
	b.ID = getter.GetString(p+"id", "")

	// b.User = getter.GetString(p+"user", "")
	// b.Driver = getter.GetString("", "")
	// b.Driver = getter.GetString("", "")

	err := getter.Error()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (inst *Demo3) loadKeys() error {

	file := fs.Default().GetPath(inst.CredentialFileName)
	text, err := file.GetIO().ReadText(nil)
	if err != nil {
		return err
	}

	src, err := collection.ParseProperties(text, nil)
	if err != nil {
		return err
	}

	dst := inst.Context.GetProperties()
	dst.Import(src.Export(nil))
	return nil
}
