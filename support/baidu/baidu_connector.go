package baidu

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type bosConnector struct {
}

func (inst *bosConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *bosConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {
	conn := &bosBucket{}
	err := conn.init(b)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
