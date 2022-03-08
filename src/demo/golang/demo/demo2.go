package demo

import (
	"strings"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// Demo2 测试大文件上传速度
type Demo2 struct {
	markup.Component `id:"demo2"`

	DemoBuckets        string              `inject:"${demo.buckets}"`
	CredentialFileName string              `inject:"${demo.credential.properties}"`
	Context            application.Context `inject:"context"`
	BM                 buckets.Manager     `inject:"#buckets.Manager"`

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

	return inst.load()
}

func (inst *Demo2) initExampleData() error {

	tmp := core.GetTempFileManager().NewTempFile()

	inst.exampleFileLocal = tmp
	inst.exampleFileName = "test/demo/object1-" + tmp.GetPath().Name()

	mk := ExampleFileMaker{}
	mk.Init(1024*1024*32, "hello")
	return mk.Make(tmp.GetPath())
}

func (inst *Demo2) load() error {

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

	driver, err := inst.BM.FindDriver(b.Driver)
	if err != nil {
		return err
	}

	p := inst.Context.GetProperties()
	b, err = driver.GetBucket("bucket", b.ID, p)
	if err != nil {
		return err
	}

	conn, err := driver.GetConnector().Open(b)
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
