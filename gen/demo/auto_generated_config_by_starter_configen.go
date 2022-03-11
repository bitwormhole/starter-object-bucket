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

	// component: demo1
	cominfobuilder.Next()
	cominfobuilder.ID("demo1").Class("").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDemo1{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: com1-demo0xbcc745.DemoRunner
	cominfobuilder.Next()
	cominfobuilder.ID("com1-demo0xbcc745.DemoRunner").Class("life").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDemoRunner{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: demo3
	cominfobuilder.Next()
	cominfobuilder.ID("demo3").Class("").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDemo3{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: demo2
	cominfobuilder.Next()
	cominfobuilder.ID("demo2").Class("").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComDemo2{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}



    return nil
}

////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDemo1 : the factory of component: demo1
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
	return nil
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
		eb.Set("com", "demo1")
		eb.Set("field", "BM")
		eb.Set("type1", "?")
		eb.Set("type2", "buckets0xc61cfb.Manager")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDemoRunner : the factory of component: com1-demo0xbcc745.DemoRunner
type comFactory4pComDemoRunner struct {

    mPrototype * demo0xbcc745.DemoRunner

	
	mDemo1Selector config.InjectionSelector
	mDemo2Selector config.InjectionSelector
	mDemo3Selector config.InjectionSelector

}

func (inst * comFactory4pComDemoRunner) init() application.ComponentFactory {

	
	inst.mDemo1Selector = config.NewInjectionSelector("#demo1",nil)
	inst.mDemo2Selector = config.NewInjectionSelector("#demo2",nil)
	inst.mDemo3Selector = config.NewInjectionSelector("#demo3",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComDemoRunner) newObject() * demo0xbcc745.DemoRunner {
	return & demo0xbcc745.DemoRunner {}
}

func (inst * comFactory4pComDemoRunner) castObject(instance application.ComponentInstance) * demo0xbcc745.DemoRunner {
	return instance.Get().(*demo0xbcc745.DemoRunner)
}

func (inst * comFactory4pComDemoRunner) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComDemoRunner) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComDemoRunner) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComDemoRunner) Init(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemoRunner) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemoRunner) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.Demo1 = inst.getterForFieldDemo1Selector(context)
	obj.Demo2 = inst.getterForFieldDemo2Selector(context)
	obj.Demo3 = inst.getterForFieldDemo3Selector(context)
	return context.LastError()
}

//getterForFieldDemo1Selector
func (inst * comFactory4pComDemoRunner) getterForFieldDemo1Selector (context application.InstanceContext) *demo0xbcc745.Demo1 {

	o1 := inst.mDemo1Selector.GetOne(context)
	o2, ok := o1.(*demo0xbcc745.Demo1)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "com1-demo0xbcc745.DemoRunner")
		eb.Set("field", "Demo1")
		eb.Set("type1", "?")
		eb.Set("type2", "*demo0xbcc745.Demo1")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}

//getterForFieldDemo2Selector
func (inst * comFactory4pComDemoRunner) getterForFieldDemo2Selector (context application.InstanceContext) *demo0xbcc745.Demo2 {

	o1 := inst.mDemo2Selector.GetOne(context)
	o2, ok := o1.(*demo0xbcc745.Demo2)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "com1-demo0xbcc745.DemoRunner")
		eb.Set("field", "Demo2")
		eb.Set("type1", "?")
		eb.Set("type2", "*demo0xbcc745.Demo2")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}

//getterForFieldDemo3Selector
func (inst * comFactory4pComDemoRunner) getterForFieldDemo3Selector (context application.InstanceContext) *demo0xbcc745.Demo3 {

	o1 := inst.mDemo3Selector.GetOne(context)
	o2, ok := o1.(*demo0xbcc745.Demo3)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "com1-demo0xbcc745.DemoRunner")
		eb.Set("field", "Demo3")
		eb.Set("type1", "?")
		eb.Set("type2", "*demo0xbcc745.Demo3")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDemo3 : the factory of component: demo3
type comFactory4pComDemo3 struct {

    mPrototype * demo0xbcc745.Demo3

	
	mDemoBucketsSelector config.InjectionSelector
	mCredentialFileNameSelector config.InjectionSelector
	mContextSelector config.InjectionSelector
	mBMSelector config.InjectionSelector

}

func (inst * comFactory4pComDemo3) init() application.ComponentFactory {

	
	inst.mDemoBucketsSelector = config.NewInjectionSelector("${demo.buckets}",nil)
	inst.mCredentialFileNameSelector = config.NewInjectionSelector("${demo.credential.properties}",nil)
	inst.mContextSelector = config.NewInjectionSelector("context",nil)
	inst.mBMSelector = config.NewInjectionSelector("#buckets.Manager",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComDemo3) newObject() * demo0xbcc745.Demo3 {
	return & demo0xbcc745.Demo3 {}
}

func (inst * comFactory4pComDemo3) castObject(instance application.ComponentInstance) * demo0xbcc745.Demo3 {
	return instance.Get().(*demo0xbcc745.Demo3)
}

func (inst * comFactory4pComDemo3) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComDemo3) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComDemo3) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComDemo3) Init(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemo3) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemo3) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.DemoBuckets = inst.getterForFieldDemoBucketsSelector(context)
	obj.CredentialFileName = inst.getterForFieldCredentialFileNameSelector(context)
	obj.Context = inst.getterForFieldContextSelector(context)
	obj.BM = inst.getterForFieldBMSelector(context)
	return context.LastError()
}

//getterForFieldDemoBucketsSelector
func (inst * comFactory4pComDemo3) getterForFieldDemoBucketsSelector (context application.InstanceContext) string {
    return inst.mDemoBucketsSelector.GetString(context)
}

//getterForFieldCredentialFileNameSelector
func (inst * comFactory4pComDemo3) getterForFieldCredentialFileNameSelector (context application.InstanceContext) string {
    return inst.mCredentialFileNameSelector.GetString(context)
}

//getterForFieldContextSelector
func (inst * comFactory4pComDemo3) getterForFieldContextSelector (context application.InstanceContext) application.Context {
    return context.Context()
}

//getterForFieldBMSelector
func (inst * comFactory4pComDemo3) getterForFieldBMSelector (context application.InstanceContext) buckets0xc61cfb.Manager {

	o1 := inst.mBMSelector.GetOne(context)
	o2, ok := o1.(buckets0xc61cfb.Manager)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "demo3")
		eb.Set("field", "BM")
		eb.Set("type1", "?")
		eb.Set("type2", "buckets0xc61cfb.Manager")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComDemo2 : the factory of component: demo2
type comFactory4pComDemo2 struct {

    mPrototype * demo0xbcc745.Demo2

	
	mDemoBucketsSelector config.InjectionSelector
	mCredentialFileNameSelector config.InjectionSelector
	mContextSelector config.InjectionSelector
	mBMSelector config.InjectionSelector

}

func (inst * comFactory4pComDemo2) init() application.ComponentFactory {

	
	inst.mDemoBucketsSelector = config.NewInjectionSelector("${demo.buckets}",nil)
	inst.mCredentialFileNameSelector = config.NewInjectionSelector("${demo.credential.properties}",nil)
	inst.mContextSelector = config.NewInjectionSelector("context",nil)
	inst.mBMSelector = config.NewInjectionSelector("#buckets.Manager",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComDemo2) newObject() * demo0xbcc745.Demo2 {
	return & demo0xbcc745.Demo2 {}
}

func (inst * comFactory4pComDemo2) castObject(instance application.ComponentInstance) * demo0xbcc745.Demo2 {
	return instance.Get().(*demo0xbcc745.Demo2)
}

func (inst * comFactory4pComDemo2) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComDemo2) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComDemo2) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComDemo2) Init(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemo2) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComDemo2) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.DemoBuckets = inst.getterForFieldDemoBucketsSelector(context)
	obj.CredentialFileName = inst.getterForFieldCredentialFileNameSelector(context)
	obj.Context = inst.getterForFieldContextSelector(context)
	obj.BM = inst.getterForFieldBMSelector(context)
	return context.LastError()
}

//getterForFieldDemoBucketsSelector
func (inst * comFactory4pComDemo2) getterForFieldDemoBucketsSelector (context application.InstanceContext) string {
    return inst.mDemoBucketsSelector.GetString(context)
}

//getterForFieldCredentialFileNameSelector
func (inst * comFactory4pComDemo2) getterForFieldCredentialFileNameSelector (context application.InstanceContext) string {
    return inst.mCredentialFileNameSelector.GetString(context)
}

//getterForFieldContextSelector
func (inst * comFactory4pComDemo2) getterForFieldContextSelector (context application.InstanceContext) application.Context {
    return context.Context()
}

//getterForFieldBMSelector
func (inst * comFactory4pComDemo2) getterForFieldBMSelector (context application.InstanceContext) buckets0xc61cfb.Manager {

	o1 := inst.mBMSelector.GetOne(context)
	o2, ok := o1.(buckets0xc61cfb.Manager)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "demo2")
		eb.Set("field", "BM")
		eb.Set("type1", "?")
		eb.Set("type2", "buckets0xc61cfb.Manager")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}




