package huawei

import (
	"errors"

	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type obsConnector struct {
}

func (inst *obsConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *obsConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {

	return nil, errors.New("no impl")
}
