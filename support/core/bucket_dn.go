package core

// BucketDNSet 这个结构保存，并管理一个存储桶相关的各个域名
type BucketDNSet struct {
	dns map[string]string
}

func (inst *BucketDNSet) getTable() map[string]string {
	table := inst.dns
	if table == nil {
		table = make(map[string]string)
		inst.dns = table
	}
	return table
}

// Init ...
func (inst *BucketDNSet) Init(kvs map[string]string) {
	dst := inst.getTable()
	for k, v := range kvs {
		dst[k] = v
	}
}

// // GetDN ...
// func (inst *BucketDNSet) GetDN(profile buckets.Profile) (string, error) {
// 	key := profile.String()
// 	table := inst.getTable()
// 	dn := table[key]
// 	if dn == "" {
// 		return "", errors.New("no domain name, with key:" + key)
// 	}
// 	return dn, nil
// }

// // SetDN ...
// func (inst *BucketDNSet) SetDN(profile buckets.Profile, dn string) {
// 	key := profile.String()
// 	table := inst.getTable()
// 	table[key] = dn
// }
