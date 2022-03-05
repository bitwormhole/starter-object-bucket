// (todo:gen2.template) 
// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package demo

import (
	buckets0xc61cfb "github.com/bitwormhole/starter-object-bucket/buckets"
	demo0xbcc745 "github.com/bitwormhole/starter-object-bucket/src/demo/golang/demo"
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

	// component: com0-demo0xbcc745.Demo1
	cominfobuilder.Next()
	cominfobuilder.ID("com0-demo0xbcc745.Demo1").Class("life").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDemo1{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}



    return nil
}

////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDemo1 : the factory of component: com0-demo0xbcc745.Demo1
type comFactory4pComDemo1 struct {

    mPrototype * demo0xbcc745.Demo1

	
	mDemoBucketsSelector config.InjectionSelector
	mCredentialFileNameSelector config.InjectionSelector
	mContextSelector config.InjectionSelector
	mBMSelector config.InjectionSelector

}

func (inst * comFactory4pComDemo1) init() application.ComponentFactory {

	
	inst.mDemoBucketsSelector = config.NewInjectionSelector("${demo.buckets}",nil)
	inst.mCredentialFileNameSelector = config.NewInjectionSelector("${demo.credential.properties}",nil)
	inst.mContextSelector = config.NewInjectionSelector("context",nil)
	inst.mBMSelector = config.NewInjectionSelector("#buckets.Manager",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComDemo1) newObject() * demo0xbcc745.Demo1 {
	return & demo0xbcc745.Demo1 {}
}

func (inst * comFactory4pComDemo1) castObject(instance application.ComponentInstance) * demo0xbcc745.Demo1 {
	return instance.Get().(*demo0xbcc745.Demo1)
}

func (inst * comFactory4pComDemo1) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComDemo1) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComDemo1) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComDemo1) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Init()
}

func (inst * comFactory4pComDemo1) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemo1) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.DemoBuckets = inst.getterForFieldDemoBucketsSelector(context)
	obj.CredentialFileName = inst.getterForFieldCredentialFileNameSelector(context)
	obj.Context = inst.getterForFieldContextSelector(context)
	obj.BM = inst.getterForFieldBMSelector(context)
	return context.LastError()
}

//getterForFieldDemoBucketsSelector
func (inst * comFactory4pComDemo1) getterForFieldDemoBucketsSelector (context application.InstanceContext) string {
    return inst.mDemoBucketsSelector.GetString(context)
}

//getterForFieldCredentialFileNameSelector
func (inst * comFactory4pComDemo1) getterForFieldCredentialFileNameSelector (context application.InstanceContext) string {
    return inst.mCredentialFileNameSelector.GetString(context)
}

//getterForFieldContextSelector
func (inst * comFactory4pComDemo1) getterForFieldContextSelector (context application.InstanceContext) application.Context {
    return context.Context()
}

//getterForFieldBMSelector
func (inst * comFactory4pComDemo1) getterForFieldBMSelector (context application.InstanceContext) buckets0xc61cfb.Manager {

	o1 := inst.mBMSelector.GetOne(context)
	o2, ok := o1.(buckets0xc61cfb.Manager)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "com0-demo0xbcc745.Demo1")
		eb.Set("field", "BM")
		eb.Set("type1", "?")
		eb.Set("type2", "buckets0xc61cfb.Manager")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}




