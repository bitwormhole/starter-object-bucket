package demo

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
)

// BucketLoader ...
type BucketLoader struct {
	markup.Component `id:"BucketLoader" initMethod:"Init"`

	CredentialFileName string              `inject:"${demo.credential.properties}"`
	Context            application.Context `inject:"context"`
	BM                 buckets.Manager     `inject:"#buckets.Manager"`
}

// Init ...
func (inst *BucketLoader) Init() error {
	return inst.loadCredentialProperties()
}

func (inst *BucketLoader) loadCredentialProperties() error {
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

// GetBucket ...
func (inst *BucketLoader) GetBucket(name string) (*buckets.Bucket, error) {
	p := inst.Context.GetProperties()
	return inst.BM.GetBucket("", name, p)
}

// OpenBucket ...
func (inst *BucketLoader) OpenBucket(name string) (buckets.Connection, error) {
	b, err := inst.GetBucket(name)
	if err != nil {
		return nil, err
	}
	return inst.BM.OpenBucket(b)
}
