package core

import (
	"errors"
	"strings"
	"sync"

	"github.com/bitwormhole/starter-object-bucket/buckets"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
)

// DefaultBucketDriverManager 是默认的存储桶驱动管理器
type DefaultBucketDriverManager struct {
	markup.Component `id:"buckets.Manager"`

	DriverSources []buckets.DriverRegistry `inject:".buckets.Driver"`

	drivers []*buckets.DriverRegistration
	mutex   sync.RWMutex
}

func (inst *DefaultBucketDriverManager) _Impl() buckets.Manager {
	return inst
}

func (inst *DefaultBucketDriverManager) getAll() []*buckets.DriverRegistration {
	list := inst.drivers
	if list == nil {
		list = inst.loadAll()
		inst.drivers = list
	}
	return list
}

func (inst *DefaultBucketDriverManager) loadAll() []*buckets.DriverRegistration {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	dst := inst.drivers
	if dst != nil {
		return dst
	}
	src := inst.DriverSources
	dst = make([]*buckets.DriverRegistration, 0)
	for _, item1 := range src {
		mid := item1.ListDrivers()
		for _, item2 := range mid {
			if item2.Driver == nil {
				continue
			}
			dst = append(dst, item2)
		}
	}
	inst.drivers = dst
	return dst
}

func (inst *DefaultBucketDriverManager) normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	return name
}

func (inst *DefaultBucketDriverManager) isNameMatch(name1, name2 string) bool {
	name1 = inst.normalizeName(name1)
	name2 = inst.normalizeName(name2)
	return name1 == name2
}

// FindDriver 查找驱动
func (inst *DefaultBucketDriverManager) FindDriver(name string) (buckets.Driver, error) {
	all := inst.getAll()
	for _, item := range all {
		if item.Driver == nil {
			continue
		}
		if inst.isNameMatch(name, item.Name) {
			return item.Driver, nil
		}
	}
	return nil, errors.New("no bucket driver with name: " + name)
}

// OpenBucket ...
func (inst *DefaultBucketDriverManager) OpenBucket(b *buckets.Bucket) (buckets.Connection, error) {
	driver, err := inst.FindDriver(b.Provider)
	if err != nil {
		return nil, err
	}
	return driver.GetConnector().Open(b)
}

// GetBucket ...
func (inst *DefaultBucketDriverManager) GetBucket(tag, id string, p collection.Properties) (*buckets.Bucket, error) {

	if tag == "" {
		tag = "bucket"
	}

	if id == "" {
		id = "default"
	}

	driverKey := tag + "." + id + ".driver"
	driverValue, err := p.GetPropertyRequired(driverKey)
	if err != nil {
		return nil, err
	}

	driver, err := inst.FindDriver(driverValue)
	if err != nil {
		return nil, err
	}

	return driver.GetBucket(tag, id, p)
}

// // ComputeDomainName ...
// func (inst *DefaultBucketDriverManager) ComputeDomainName(b *buckets.Bucket, profile string) (string, error) {
// 	driver := b.Driver
// 	if driver == nil {
// 		drv, err := inst.FindDriver(b.Provider)
// 		if err != nil {
// 			return "", err
// 		}
// 		driver = drv
// 	}
// }
