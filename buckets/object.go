package buckets

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/util"
)

// ObjectMeta 表示对象的元数据
type ObjectMeta struct {
	Size          int64
	Date          time.Time
	Hash          string
	HashAlgorithm string
	ContentType   string
	More          map[string]string
}

// HTTPUploading ...
type HTTPUploading struct {
	UseHTTPS       bool
	Method         string
	URL            string
	ContentType    string
	ContentLength  int64
	ContentMD5     util.Hex
	MaxAge         int64 // in seconds
	RequestHeaders map[string]string
}

// ObjectEntity 表示对象的实体
type ObjectEntity interface {
	GetSize() int64
	Open() (io.ReadCloser, error)
}

// Object 表示对一个对象的引用
type Object interface {
	Exists() (bool, error)
	GetDownloadURL() string
	GetMeta() (*ObjectMeta, error)
	GetName() string
	UpdateMeta(meta *ObjectMeta) error
	GetEntity() (ObjectEntity, error)
	PutEntity(entity ObjectEntity, meta *ObjectMeta) error
	PutFile(file fs.Path, meta *ObjectMeta) error
	// GetFile(saveAs fs.Path) error

	// UploadByAPI 准备上传参数;
	// 输入：待准备的参数;
	// 输出：已经准备好的参数;
	UploadByAPI(up *HTTPUploading) (*HTTPUploading, error)
}

////////////////////////////////////////////////////////////////////////////////

// NewEntityForFile ...
func NewEntityForFile(file fs.Path) (ObjectEntity, error) {
	if !file.IsFile() {
		return nil, errors.New("the file is not exists, path=" + file.Path())
	}
	ent := &fileEntity{}
	ent.file = file
	return ent, nil
}

// NewEntityForFileName ...
func NewEntityForFileName(filename string) (ObjectEntity, error) {
	file := fs.Default().GetPath(filename)
	return NewEntityForFile(file)
}

type fileEntity struct {
	file fs.Path
}

func (inst *fileEntity) GetSize() int64 {
	return inst.file.Size()
}

func (inst *fileEntity) Open() (io.ReadCloser, error) {
	return inst.file.GetIO().OpenReader(nil)
}

////////////////////////////////////////////////////////////////////////////////

// NewEntityForBytes ...
func NewEntityForBytes(b []byte) (ObjectEntity, error) {
	ent := &ramEntity{}
	ent.data = b
	return ent, nil
}

type ramEntity struct {
	data []byte
}

func (inst *ramEntity) GetSize() int64 {
	n := len(inst.data)
	return int64(n)
}

func (inst *ramEntity) Open() (io.ReadCloser, error) {
	r1 := bytes.NewReader(inst.data)
	r2 := io.NopCloser(r1)
	return r2, nil
}

////////////////////////////////////////////////////////////////////////////////
