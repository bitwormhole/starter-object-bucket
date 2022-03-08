package huawei

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type obsConnector struct {
}

func (inst *obsConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *obsConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {
	conn := &obsBucket{}
	err := conn.init(b)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
