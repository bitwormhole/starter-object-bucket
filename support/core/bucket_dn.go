package core

import (
	"errors"

	"github.com/bitwormhole/starter-object-bucket/buckets"
)

// BucketDNSet 这个结构保存，并管理一个存储桶相关的各个域名
type BucketDNSet struct {
	dns map[buckets.DomainType]string
}

func (inst *BucketDNSet) getTable() map[buckets.DomainType]string {
	table := inst.dns
	if table == nil {
		table = make(map[buckets.DomainType]string)
		inst.dns = table
	}
	return table
}

func (inst *BucketDNSet) Init(kvs map[string]string) {
	dst := inst.getTable()
	for k, v := range kvs {
		dntype := buckets.DomainType(k)
		dst[dntype] = v
	}
}

func (inst *BucketDNSet) GetDN(dntype buckets.DomainType) (string, error) {
	table := inst.getTable()
	dn := table[dntype]
	if dn == "" {
		key := string(dntype)
		return "", errors.New("no domain name, with key:" + key)
	}
	return dn, nil
}

func (inst *BucketDNSet) SetDN(dntype buckets.DomainType, dn string) {
	table := inst.getTable()
	table[dntype] = dn
}
