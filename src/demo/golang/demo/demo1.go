package demo

import (
	"encoding/json"
	"strings"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter-object-bucket/support/core"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// Demo1 测试功能是否正常
type Demo1 struct {
	markup.Component `id:"demo1"`

	DemoBuckets        string              `inject:"${demo.buckets}"`
	CredentialFileName string              `inject:"${demo.credential.properties}"`
	Context            application.Context `inject:"context"`
	BM                 buckets.Manager     `inject:"#buckets.Manager"`

	exampleFileName  string
	exampleFileLocal core.TempFile
}

func (inst *Demo1) _Impl() application.LifeRegistry {
	return inst
}

// GetLifeRegistration ...
func (inst *Demo1) GetLifeRegistration() *application.LifeRegistration {
	lr := &application.LifeRegistration{}
	lr.OnInit = inst.Init
	lr.Looper = inst
	return lr
}

// Init ...
func (inst *Demo1) Init() error {

	err := inst.initExampleData()
	if err != nil {
		return err
	}

	return inst.load()
}

func (inst *Demo1) initExampleData() error {

	tmp := core.GetTempFileManager().NewTempFile()

	inst.exampleFileLocal = tmp
	inst.exampleFileName = "test/demo/object1-" + tmp.GetPath().Name()

	mk := ExampleFileMaker{}
	mk.Init(1024*32, "hello")
	return mk.Make(tmp.GetPath())
}

func (inst *Demo1) load() error {

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
func (inst *Demo1) Loop() error {

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

func (inst *Demo1) testBucket(b *buckets.Bucket) error {

	driver, err := inst.BM.FindDriver(b.Provider)
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

	return inst.testObject(o)
}

func (inst *Demo1) testObject(o buckets.Object) error {

	const tab = "    "

	// Exists
	ext, err := o.Exists()
	if err != nil {
		return err
	}
	vlog.Info(tab+"o.exists=", ext)

	// GetMeta
	meta, err := o.GetMeta()
	if err != nil {
		return err
	}
	vlog.Info(tab+"o.meta=", inst.stringifyMeta(meta))

	// GetDownloadURL
	url := o.GetDownloadURL()
	vlog.Info(tab+"o.dl_url=", url)

	// GetDownloadURL
	name := o.GetName()
	vlog.Info(tab+"o.name=", name)

	// o.GetEntity()

	// UpdateMeta
	err = o.UpdateMeta(meta)
	if err != nil {
		vlog.Warn("UpdateMeta:error:", err)
		// return err
	}

	return nil
}

func (inst *Demo1) stringifyMeta(meta *buckets.ObjectMeta) string {
	j, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		return "[error]"
	}
	return string(j)
}

func (inst *Demo1) loadBuckets() ([]*buckets.Bucket, error) {
	dst := make([]*buckets.Bucket, 0)
	items := strings.Split(inst.DemoBuckets, ",")
	for _, item := range items {
		b, err := inst.loadBucketWithID(item)
		if err != nil {
			return nil, err
		}
		dst = append(dst, b)
	}
	return dst, nil
}

func (inst *Demo1) loadBucketWithID(id string) (*buckets.Bucket, error) {

	id = strings.TrimSpace(id)

	b := &buckets.Bucket{}
	p := "bucket." + id + "."
	getter := inst.Context.GetProperties().Getter()

	b.Provider = getter.GetString(p+"driver", "")
	b.ID = id

	// b.User = getter.GetString(p+"user", "")
	// b.Driver = getter.GetString("", "")
	// b.Driver = getter.GetString("", "")

	err := getter.Error()
	if err != nil {
		return nil, err
	}
	return b, nil
}
