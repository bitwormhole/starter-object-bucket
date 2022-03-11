// (todo:gen2.template) 
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
	application "github.com/bitwormhole/starter/application"
	config "github.com/bitwormhole/starter/application/config"
	lang "github.com/bitwormhole/starter/lang"
	util "github.com/bitwormhole/starter/util"
    
)


func nop(x ... interface{}){
	util.Int64ToTime(0)
	lang.CreateReleasePool()
}


func autoGenConfig(cb application.ConfigBuilder) error {

	var err error = nil
	cominfobuilder := config.ComInfo()
	nop(err,cominfobuilder)

	// component: com0-huawei0xcc0ad2.OBSDriver
	cominfobuilder.Next()
	cominfobuilder.ID("com0-huawei0xcc0ad2.OBSDriver").Class("buckets.Driver").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComOBSDriver{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: com1-qq0x13bfdf.COSDriver
	cominfobuilder.Next()
	cominfobuilder.ID("com1-qq0x13bfdf.COSDriver").Class("buckets.Driver").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComCOSDriver{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: com2-aliyun0x8a6680.OSSDriver
	cominfobuilder.Next()
	cominfobuilder.ID("com2-aliyun0x8a6680.OSSDriver").Class("buckets.Driver").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComOSSDriver{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: com3-baidu0xa9ed2f.BOSDriver
	cominfobuilder.Next()
	cominfobuilder.ID("com3-baidu0xa9ed2f.BOSDriver").Class("buckets.Driver").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComBOSDriver{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: buckets.Manager
	cominfobuilder.Next()
	cominfobuilder.ID("buckets.Manager").Class("").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDefaultBucketDriverManager{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}



    return nil
}

////////////////////////////////////////////////////////////////////////////////

// comFactory4pComOBSDriver : the factory of component: com0-huawei0xcc0ad2.OBSDriver
type comFactory4pComOBSDriver struct {

    mPrototype * huawei0xcc0ad2.OBSDriver

	

}

func (inst * comFactory4pComOBSDriver) init() application.ComponentFactory {

	


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComOBSDriver) newObject() * huawei0xcc0ad2.OBSDriver {
	return & huawei0xcc0ad2.OBSDriver {}
}

func (inst * comFactory4pComOBSDriver) castObject(instance application.ComponentInstance) * huawei0xcc0ad2.OBSDriver {
	return instance.Get().(*huawei0xcc0ad2.OBSDriver)
}

func (inst * comFactory4pComOBSDriver) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComOBSDriver) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComOBSDriver) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComOBSDriver) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Init()
}

func (inst * comFactory4pComOBSDriver) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComOBSDriver) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	return nil
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComCOSDriver : the factory of component: com1-qq0x13bfdf.COSDriver
type comFactory4pComCOSDriver struct {

    mPrototype * qq0x13bfdf.COSDriver

	

}

func (inst * comFactory4pComCOSDriver) init() application.ComponentFactory {

	


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComCOSDriver) newObject() * qq0x13bfdf.COSDriver {
	return & qq0x13bfdf.COSDriver {}
}

func (inst * comFactory4pComCOSDriver) castObject(instance application.ComponentInstance) * qq0x13bfdf.COSDriver {
	return instance.Get().(*qq0x13bfdf.COSDriver)
}

func (inst * comFactory4pComCOSDriver) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComCOSDriver) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComCOSDriver) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComCOSDriver) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Init()
}

func (inst * comFactory4pComCOSDriver) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComCOSDriver) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	return nil
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComOSSDriver : the factory of component: com2-aliyun0x8a6680.OSSDriver
type comFactory4pComOSSDriver struct {

    mPrototype * aliyun0x8a6680.OSSDriver

	

}

func (inst * comFactory4pComOSSDriver) init() application.ComponentFactory {

	


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComOSSDriver) newObject() * aliyun0x8a6680.OSSDriver {
	return & aliyun0x8a6680.OSSDriver {}
}

func (inst * comFactory4pComOSSDriver) castObject(instance application.ComponentInstance) * aliyun0x8a6680.OSSDriver {
	return instance.Get().(*aliyun0x8a6680.OSSDriver)
}

func (inst * comFactory4pComOSSDriver) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComOSSDriver) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComOSSDriver) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComOSSDriver) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Init()
}

func (inst * comFactory4pComOSSDriver) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComOSSDriver) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	return nil
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComBOSDriver : the factory of component: com3-baidu0xa9ed2f.BOSDriver
type comFactory4pComBOSDriver struct {

    mPrototype * baidu0xa9ed2f.BOSDriver

	

}

func (inst * comFactory4pComBOSDriver) init() application.ComponentFactory {

	


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComBOSDriver) newObject() * baidu0xa9ed2f.BOSDriver {
	return & baidu0xa9ed2f.BOSDriver {}
}

func (inst * comFactory4pComBOSDriver) castObject(instance application.ComponentInstance) * baidu0xa9ed2f.BOSDriver {
	return instance.Get().(*baidu0xa9ed2f.BOSDriver)
}

func (inst * comFactory4pComBOSDriver) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComBOSDriver) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComBOSDriver) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComBOSDriver) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Init()
}

func (inst * comFactory4pComBOSDriver) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComBOSDriver) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	return nil
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDefaultBucketDriverManager : the factory of component: buckets.Manager
type comFactory4pComDefaultBucketDriverManager struct {

    mPrototype * core0xd5361e.DefaultBucketDriverManager

	
	mDriverSourcesSelector config.InjectionSelector

}

func (inst * comFactory4pComDefaultBucketDriverManager) init() application.ComponentFactory {

	
	inst.mDriverSourcesSelector = config.NewInjectionSelector(".buckets.Driver",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComDefaultBucketDriverManager) newObject() * core0xd5361e.DefaultBucketDriverManager {
	return & core0xd5361e.DefaultBucketDriverManager {}
}

func (inst * comFactory4pComDefaultBucketDriverManager) castObject(instance application.ComponentInstance) * core0xd5361e.DefaultBucketDriverManager {
	return instance.Get().(*core0xd5361e.DefaultBucketDriverManager)
}

func (inst * comFactory4pComDefaultBucketDriverManager) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComDefaultBucketDriverManager) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComDefaultBucketDriverManager) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComDefaultBucketDriverManager) Init(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDefaultBucketDriverManager) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDefaultBucketDriverManager) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.DriverSources = inst.getterForFieldDriverSourcesSelector(context)
	return context.LastError()
}

//getterForFieldDriverSourcesSelector
func (inst * comFactory4pComDefaultBucketDriverManager) getterForFieldDriverSourcesSelector (context application.InstanceContext) []buckets0xc61cfb.DriverRegistry {
	list1 := inst.mDriverSourcesSelector.GetList(context)
	list2 := make([]buckets0xc61cfb.DriverRegistry, 0, len(list1))
	for _, item1 := range list1 {
		item2, ok := item1.(buckets0xc61cfb.DriverRegistry)
		if ok {
			list2 = append(list2, item2)
		}
	}
	return list2
}




