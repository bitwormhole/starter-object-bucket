package demo

import (
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/markup"
)

// DemoRunner ...
type DemoRunner struct {
	markup.Component `class:"life"`

	DemoSelector string `inject:"${demo.selector}"`
	Demo1        *Demo1 `inject:"#demo1"`
	Demo2        *Demo2 `inject:"#demo2"`
	Demo3        *Demo3 `inject:"#demo3"`

	current application.LifeRegistration
}

func (inst *DemoRunner) _Impl() application.LifeRegistry {
	return inst
}

// GetLifeRegistration ...
func (inst *DemoRunner) GetLifeRegistration() *application.LifeRegistration {
	lr := &application.LifeRegistration{}
	lr.OnInit = inst.Init
	lr.OnStart = inst.Start
	lr.Looper = inst
	lr.OnStop = inst.Stop
	lr.OnDestroy = inst.Destroy
	return lr
}

func (inst *DemoRunner) selectCurrent() application.LifeRegistry {
	name := inst.DemoSelector
	all := map[string]application.LifeRegistry{
		"demo1": inst.Demo1,
		"demo2": inst.Demo2,
		"demo3": inst.Demo3,
	}
	lr := all[name]
	if lr == nil {
		panic("no demo with name:" + name)
	}
	return lr
}

func (inst *DemoRunner) initCurrent() {

	lr := inst.selectCurrent().GetLifeRegistration()
	nop := &demoRunnerNOP{}

	if lr.OnInit == nil {
		lr.OnInit = nop.Loop
	}
	if lr.OnStart == nil {
		lr.OnStart = nop.Loop
	}
	if lr.Looper == nil {
		lr.Looper = nop
	}
	if lr.OnStop == nil {
		lr.OnStop = nop.Loop
	}
	if lr.OnDestroy == nil {
		lr.OnDestroy = nop.Loop
	}

	inst.current = *lr
}

// Init ...
func (inst *DemoRunner) Init() error {

	inst.initCurrent()

	return inst.current.OnInit()
}

// Start ...
func (inst *DemoRunner) Start() error {
	return inst.current.OnStart()
}

// Loop ...
func (inst *DemoRunner) Loop() error {
	return inst.current.Looper.Loop()
}

// Stop ...
func (inst *DemoRunner) Stop() error {
	return inst.current.OnStop()
}

// Destroy ...
func (inst *DemoRunner) Destroy() error {
	return inst.current.OnDestroy()
}

////////////////////////////////////////////////////////////////////////////////

type demoRunnerNOP struct {
}

func (inst *demoRunnerNOP) Loop() error {
	return nil // NOP
}
