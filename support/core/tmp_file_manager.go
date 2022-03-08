package core

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/util"
)

// TempFile 代表临时文件
type TempFile interface {
	io.Closer
	GetPath() fs.Path
}

// TempFileManager 代表临时文件管理器
type TempFileManager interface {
	NewTempFile() TempFile
	GetTempDir() fs.Path
	SetTempDir(dir fs.Path)
}

////////////////////////////////////////////////////////////////////////////////

type innerTempFileManager struct {
	mutex      sync.Mutex
	tmpdir     fs.Path
	initialled bool

	index          int64
	startupTime    util.Time
	fileNamePrefix string
}

func (inst *innerTempFileManager) _Impl() TempFileManager {
	return inst
}

func (inst *innerTempFileManager) init() {

	if inst.initialled {
		return
	}

	now := util.Now()
	nonce := make([]byte, 10)
	builder := bytes.Buffer{}

	rand.Reader.Read(nonce)
	builder.Write(nonce)
	builder.WriteString(now.String())
	sum := md5.Sum(builder.Bytes())

	inst.index = 1
	inst.startupTime = now
	inst.fileNamePrefix = util.StringifyBytes(sum[:])
	inst.initialled = true
}

func (inst *innerTempFileManager) NewTempFile() TempFile {

	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	inst.index++
	index := inst.index
	now := util.Now()
	dir := inst.GetTempDir()

	builder := strings.Builder{}
	builder.WriteString(inst.fileNamePrefix)
	builder.WriteString("_")
	builder.WriteString(strconv.FormatInt(now.Int64(), 10))
	builder.WriteString("_")
	builder.WriteString(strconv.FormatInt(index, 10))
	builder.WriteString(".tmp~")

	if !dir.Exists() {
		dir.Mkdirs()
	}
	file := dir.GetChild(builder.String())
	tmp := &innerTempFile{file: file}
	return tmp
}

func (inst *innerTempFileManager) loadTempDir() fs.Path {
	dir := os.TempDir()
	return fs.Default().GetPath(dir + "/starter-object-bucket")
}

func (inst *innerTempFileManager) GetTempDir() fs.Path {
	dir := inst.tmpdir
	if dir == nil {
		dir = inst.loadTempDir()
		inst.tmpdir = dir
	}
	return dir
}

func (inst *innerTempFileManager) SetTempDir(dir fs.Path) {
	inst.tmpdir = dir
}

////////////////////////////////////////////////////////////////////////////////

type innerTempFile struct {
	file fs.Path
}

func (inst *innerTempFile) _Impl() TempFile {
	return inst
}

func (inst *innerTempFile) GetPath() fs.Path {
	return inst.file
}

func (inst *innerTempFile) Close() error {
	file := inst.file
	inst.file = nil
	if file == nil {
		return nil
	}
	if file.Exists() && file.IsFile() {
		file.Delete()
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

var theTempFileManager innerTempFileManager

// GetTempFileManager 取临时文件管理器
func GetTempFileManager() TempFileManager {
	man := &theTempFileManager
	man.init()
	return man
}

////////////////////////////////////////////////////////////////////////////////

type tmpFileHolder struct {
	tmp TempFile
}

func (inst *tmpFileHolder) close() {
	t := inst.tmp
	inst.tmp = nil
	if t == nil {
		return
	}
	t.Close()
}

func (inst *tmpFileHolder) clear() {
	inst.tmp = nil
}

// PrepareLargeTempFileForUploading 为上传大型文件做准备
func PrepareLargeTempFileForUploading(entity buckets.ObjectEntity) (TempFile, error) {

	tmp := GetTempFileManager().NewTempFile()
	holder := tmpFileHolder{tmp: tmp}
	defer holder.close()
	file := tmp.GetPath()
	opt := file.FileSystem().DefaultWriteOptions()
	opt.Create = true

	dst, err := file.GetIO().OpenWriter(opt, true)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	src, err := entity.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buffer := make([]byte, 64*1024)
	err = util.PumpStream(src, dst, buffer)
	if err != nil {
		return nil, err
	}

	holder.clear()
	return tmp, nil
}

////////////////////////////////////////////////////////////////////////////////
