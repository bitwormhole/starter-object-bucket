package demo

import (
	"strings"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// Demo1 ...
type Demo1 struct {
	markup.Component `class:"life" initMethod:"Init"`

	DemoBuckets        string              `inject:"${demo.buckets}"`
	CredentialFileName string              `inject:"${demo.credential.properties}"`
	Context            application.Context `inject:"context"`
	BM                 buckets.Manager     `inject:"#buckets.Manager"`
}

func (inst *Demo1) _Impl() application.LifeRegistry {
	return inst
}

// Init ...
func (inst *Demo1) Init() error {
	return inst.load()
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

// GetLifeRegistration ...
func (inst *Demo1) GetLifeRegistration() *application.LifeRegistration {
	lr := &application.LifeRegistration{}
	lr.Looper = inst
	return lr
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

	data := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz+/"
	entity, err := buckets.NewEntityForBytes([]byte(data))
	if err != nil {
		return err
	}

	o := conn.GetObject("test/demo/object1")
	err = o.PutEntity(entity, nil)
	if err != nil {
		return err
	}

	return nil
}

func (inst *Demo1) loadBuckets() ([]*buckets.Bucket, error) {
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

func (inst *Demo1) loadBucketWithName(name string) (*buckets.Bucket, error) {

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
