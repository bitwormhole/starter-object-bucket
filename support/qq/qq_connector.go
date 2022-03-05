package qq

import (
	"errors"

	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type cosConnector struct {
}

func (inst *cosConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *cosConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {

	return nil, errors.New("no impl")
}
