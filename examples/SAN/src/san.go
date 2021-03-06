package src

import (
	lib "gosan"
	"gosan/deepcontroller"
	"log"
)

/*
Implement Builder interface
*/

type (
	SANBuilder interface {
		IOBalancer(*IOBalancer) SANBuilder
		Controllers([]*Controller) SANBuilder
		JbodControllers([]*SANJBODController) SANBuilder
		Client(*ClientManager) SANBuilder
		Tracer(*TracerManager) SANBuilder
		AtmosphereControl(*AtmosphereControlManager) SANBuilder
		Build() SHD
	}
	sanBuilder struct {
		cm  *ClientManager
		iob *IOBalancer
		vc  []*Controller
		jc  []*SANJBODController

		tm       *TracerManager
		aControl *AtmosphereControlManager
	}
)

type (
	SAN struct {
		cm  *ClientManager
		iob *IOBalancer
		vc  []*Controller
		jc  []*SANJBODController

		tm       *TracerManager
		aControl *AtmosphereControlManager
	}
	SHD interface {
		Read()
		Write()
		Save()

		GetControllers() []*Controller
		GetClientManager() *ClientManager
		GetIoBalancer() *IOBalancer
		GetSANJBODControllers() []*SANJBODController

		GetJsonLog() *LogJson
		GetDeepControllerManager() *deepcontroller.Manager
	}
)

func (t *sanBuilder) IOBalancer(iob *IOBalancer) SANBuilder {
	t.iob = iob
	return t
}

func (t *sanBuilder) Controllers(cons []*Controller) SANBuilder {
	t.vc = cons
	return t
}

func (t *sanBuilder) JbodControllers(jc []*SANJBODController) SANBuilder {
	t.jc = jc
	return t
}

func (t *sanBuilder) Client(cm *ClientManager) SANBuilder {
	t.cm = cm
	return t
}

func (t *sanBuilder) Tracer(tm *TracerManager) SANBuilder {
	t.tm = tm
	return t
}

func (t *sanBuilder) AtmosphereControl(ac *AtmosphereControlManager) SANBuilder {
	t.aControl = ac
	return t
}

func (t *sanBuilder) Build() SHD {
	return &SAN{
		cm:       t.cm,
		iob:      t.iob,
		vc:       t.vc,
		jc:       t.jc,
		tm:       t.tm,
		aControl: t.aControl,
	}
}

func NewSAN() SANBuilder {
	return &sanBuilder{}
}
func (t *SAN) Read() {
	panic("implement me")
}

func (t *SAN) Write() {
	panic("implement me")
}

func (t *SAN) Save() {
	panic("implement me")
}

func (t *SAN) GetControllers() []*Controller {
	return t.vc
}

func (t *SAN) GetClientManager() *ClientManager {
	return t.cm
}

func (t *SAN) GetIoBalancer() *IOBalancer {
	return t.iob
}

func (t *SAN) GetSANJBODControllers() []*SANJBODController {
	return t.jc
}

func (t *SAN) GetJsonLog() *LogJson {
	return t.tm.logs
}

func (t *SAN) GetDeepControllerManager() *deepcontroller.Manager {
	panic("implement me")
}

type SANComponent struct {
	currentState string `json:"-"`
}

func (component *SANComponent) GetCurrentState() string {
	if component.currentState == "" {
		log.Panic("Current state is empty")
	}
	return component.currentState
}

func (component *SANComponent) SetCurrentState(state string) {
	if state == "" {
		log.Panic("state is empty")
	}
	component.currentState = state
}

/*

===========


===========
*/

type (
	SANBFunction func(*SANProcess, interface{})

	SANProcess struct {
		*lib.Process
	}
)

func NewSANProcess() *SANProcess {
	return &SANProcess{}
}

func (tp *SANProcess) SetParent(p *lib.Process) {
	tp.Process = p
}

func (tp *SANProcess) GetParent() *lib.Process {
	return tp.Process
}

func FORK(name string, f SANBFunction, host *lib.Host, data interface{}) {
	lib.SIM_subprocess_create_async(name, PIDecorator(f), host, NewSANProcess(), data)
}

func PIDecorator(function SANBFunction) lib.BaseBehaviourFunction {
	return func(pi lib.ProcessInterface, data interface{}) {

		tp, ok := pi.(*SANProcess)
		if !ok {
			log.Panic("SAN machine conversion error")
		}
		function(tp, data)
	}
}
