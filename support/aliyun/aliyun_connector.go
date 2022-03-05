package aliyun

import (
	"github.com/bitwormhole/starter-object-bucket/buckets"
)

type aliyunConnector struct {
}

func (inst *aliyunConnector) _Impl() buckets.Connector {
	return inst
}

func (inst *aliyunConnector) Open(b *buckets.Bucket) (buckets.Connection, error) {
	conn := &ossBucketConnection{}
	err := conn.init(b)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
