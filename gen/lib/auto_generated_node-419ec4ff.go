// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package lib

import (
	buckets0xc61cfb "github.com/bitwormhole/starter-object-bucket/buckets"
	aliyun0x8a6680 "github.com/bitwormhole/starter-object-bucket/support/aliyun"
	baidu0xa9ed2f "github.com/bitwormhole/starter-object-bucket/support/baidu"
	core0xd5361e "github.com/bitwormhole/starter-object-bucket/support/core"
	huawei0xcc0ad2 "github.com/bitwormhole/starter-object-bucket/support/huawei"
	qq0x13bfdf "github.com/bitwormhole/starter-object-bucket/support/qq"
	markup0x23084a "github.com/bitwormhole/starter/markup"
)

type pComOBSDriver struct {
	instance *huawei0xcc0ad2.OBSDriver
	 markup0x23084a.Component `class:"buckets.Driver" initMethod:"Init"`
}


type pComCOSDriver struct {
	instance *qq0x13bfdf.COSDriver
	 markup0x23084a.Component `class:"buckets.Driver" initMethod:"Init"`
}


type pComOSSDriver struct {
	instance *aliyun0x8a6680.OSSDriver
	 markup0x23084a.Component `class:"buckets.Driver" initMethod:"Init"`
}


type pComBOSDriver struct {
	instance *baidu0xa9ed2f.BOSDriver
	 markup0x23084a.Component `class:"buckets.Driver" initMethod:"Init"`
}


type pComDefaultBucketDriverManager struct {
	instance *core0xd5361e.DefaultBucketDriverManager
	 markup0x23084a.Component `id:"buckets.Manager"`
	DriverSources []buckets0xc61cfb.DriverRegistry `inject:".buckets.Driver"`
}

