package demo

import (
	"strings"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// Demo2 测试大文件上传速度
type Demo2 struct {
	markup.Component `id:"demo2"`

	// CredentialFileName string              `inject:"${demo.credential.properties}"`
	// Context            application.Context `inject:"context"`

	DemoBuckets  string          `inject:"${demo.buckets}"`
	BM           buckets.Manager `inject:"#buckets.Manager"`
	BucketLoader *BucketLoader   `inject:"#BucketLoader"`

	exampleFileName  string
	exampleFileLocal core.TempFile
}

func (inst *Demo2) _Impl() application.LifeRegistry {
	return inst
}

// GetLifeRegistration ...
func (inst *Demo2) GetLifeRegistration() *application.LifeRegistration {
	lr := &application.LifeRegistration{}
	lr.OnInit = inst.Init
	lr.Looper = inst
	return lr
}

// Init ...
func (inst *Demo2) Init() error {
	err := inst.initExampleData()
	if err != nil {
		return err
	}
	return nil
}

func (inst *Demo2) initExampleData() error {

	tmp := core.GetTempFileManager().NewTempFile()

	inst.exampleFileLocal = tmp
	inst.exampleFileName = "test/demo/object1-" + tmp.GetPath().Name()

	mk := ExampleFileMaker{}
	mk.Init(1024*1024*32, "hello")
	return mk.Make(tmp.GetPath())
}

// Loop ...
func (inst *Demo2) Loop() error {

	vlog.Warn("todo: loop")

	blist, err := inst.loadBuckets()
	if err != nil {
		return err
	}

	for _, item := range blist {
		err = inst.testBucket(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (inst *Demo2) testBucket(b *buckets.Bucket) error {

	conn, err := inst.BucketLoader.OpenBucket(b.ID)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Check()
	if err != nil {
		return err
	}

	file := inst.exampleFileLocal.GetPath()
	entity, err := buckets.NewEntityForFile(file)
	if err != nil {
		return err
	}
	fileSize := file.Size()

	vlog.Info("begin: upload to ", b.Driver)
	defer func() {
		vlog.Info("end  : file_size=", fileSize)
	}()

	o := conn.GetObject(inst.exampleFileName)
	err = o.PutEntity(entity, nil)
	if err != nil {
		return err
	}

	return nil
}

func (inst *Demo2) loadBuckets() ([]*buckets.Bucket, error) {
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

func (inst *Demo2) loadBucketWithName(name string) (*buckets.Bucket, error) {
	return inst.BucketLoader.GetBucket(name)
}
