// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package demo

import (
	buckets0xc61cfb "github.com/bitwormhole/starter-object-bucket/buckets"
	demo0xbcc745 "github.com/bitwormhole/starter-object-bucket/src/demo/golang/demo"
	application0x67f6c5 "github.com/bitwormhole/starter/application"
	markup0x23084a "github.com/bitwormhole/starter/markup"
)

type pComDemo1 struct {
	instance *demo0xbcc745.Demo1
	 markup0x23084a.Component `id:"demo1"`
	DemoBuckets string `inject:"${demo.buckets}"`
	CredentialFileName string `inject:"${demo.credential.properties}"`
	Context application0x67f6c5.Context `inject:"context"`
	BM buckets0xc61cfb.Manager `inject:"#buckets.Manager"`
}


type pComDemoRunner struct {
	instance *demo0xbcc745.DemoRunner
	 markup0x23084a.Component `class:"life"`
	Demo1 *demo0xbcc745.Demo1 `inject:"#demo1"`
	Demo2 *demo0xbcc745.Demo2 `inject:"#demo2"`
	Demo3 *demo0xbcc745.Demo3 `inject:"#demo3"`
}


type pComDemo3 struct {
	instance *demo0xbcc745.Demo3
	 markup0x23084a.Component `id:"demo3"`
	DemoBuckets string `inject:"${demo.buckets}"`
	CredentialFileName string `inject:"${demo.credential.properties}"`
	Context application0x67f6c5.Context `inject:"context"`
	BM buckets0xc61cfb.Manager `inject:"#buckets.Manager"`
}


type pComDemo2 struct {
	instance *demo0xbcc745.Demo2
	 markup0x23084a.Component `id:"demo2"`
	DemoBuckets string `inject:"${demo.buckets}"`
	CredentialFileName string `inject:"${demo.credential.properties}"`
	Context application0x67f6c5.Context `inject:"context"`
	BM buckets0xc61cfb.Manager `inject:"#buckets.Manager"`
}

