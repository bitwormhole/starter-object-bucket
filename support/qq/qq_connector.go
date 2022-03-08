package qq

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type cosConnector struct {
}

func (inst *cosConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *cosConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {
	b2 := &cosBucket{}
	err := b2.init(b)
	if err != nil {
		return nil, err
	}
	return b2, nil
}
