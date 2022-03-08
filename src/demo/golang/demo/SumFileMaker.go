package demo

import (
	"crypto/sha256"
	"errors"

	"github.com/bitwormhole/starter/io/fs"
)

type ExampleFileMaker struct {
	seed     []byte
	wantSize int64
}

func (inst *ExampleFileMaker) Init(size int64, seed string) {
	s := []byte(seed)
	inst.seed = s
	inst.wantSize = size
}

func (inst *ExampleFileMaker) Make(file fs.Path) error {

	opt := file.FileSystem().DefaultWriteOptions()
	opt.Create = true

	dst, err := file.GetIO().OpenWriter(opt, true)
	if err != nil {
		return err
	}
	defer dst.Close()
	sum := inst.seed
	want := inst.wantSize
	count := int64(0)
	for count < want {
		s2 := sha256.Sum256(sum)
		sum = s2[:]
		size1 := want - count
		size2 := int64(len(sum))
		if size1 < size2 {
			size2 = size1
		} else {
			size1 = size2
		}
		n, err := dst.Write(sum[0:size1])
		if n > 0 {
			count += int64(n)
		}
		if err != nil {
			return err
		}
		if n != int(size1) {
			return errors.New("bad size of write")
		}
	}
	return nil
}
