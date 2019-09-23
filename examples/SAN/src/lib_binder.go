package src

import (
	"lib"
	"log"
)

// Reference function
func Referencefunction(tp *SANProcess, data interface{}) {}

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

func FORK_SYNC(name string, f SANBFunction, host *lib.Host, data interface{}) {
	lib.SIM_subprocess_create_sync(name, PIDecorator(f), host, NewSANProcess(), data)
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
