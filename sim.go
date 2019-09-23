package lib

import (
	"container/heap"
	"fmt"
	"log"
)

var _ = fmt.Println

type BaseBehaviourFunction func(ProcessInterface, interface{})

func SIM_init(packet_config string) {
	env := NewEnvironment()

	SIM_get_clock = SIM_get_clockFactory(env)
	SIM_run = SIM_runFactory(env)
	GetCallbacks = initCallbacksFactory(env)
	NewProcess = NewProcessFactory(env)

	InitPacketInfo(packet_config)
}

func SIM_runFactory(env *Environment) func(float64) {
	return func(until float64) {
		if until != -1 {
			globalStop := NewEvent(until, nil, StopEvent, nil)
			heap.Push(&env.globalQueue, globalStop)
		}

		for !env.shouldStop {
			env.SendResumeSignalToWorkers()
			env.Step()
		}
		log.Printf("Simulation took %.2f\n", SIM_get_clock())

		//Close resources (file, tcp-connection)
		for _, resource := range env.systemResources {
			resource.Close()
		}
	}
}

func SIM_get_clockFactory(env *Environment) func() float64 {
	return func() float64 {
		return env.currentTime
	}
}

func SIM_subprocess_create_async(name string, f BaseBehaviourFunction, host *Host, pi ProcessInterface,
	data interface{}) ProcessInterface {

	p := NewProcess(name, host)
	pi.SetParent(p)

	ProcWrapper(f, pi, data)
	return pi
}

func SIM_subprocess_create_sync(name string, f BaseBehaviourFunction, host *Host, pi ProcessInterface, data interface{}) ProcessInterface {
	pi = SIM_subprocess_create_async(name, f, host, pi, data)

	pi.GetParent().env.stepEnd <- struct{}{}
	<-pi.GetParent().resumeChan
	return pi
}

func sim_function_registerFactory(FunctionsMap map[string]BaseBehaviourFunction, FunctionsData map[string]interface{}) func(string, BaseBehaviourFunction, interface{}) {
	return func(FuncName string, Func BaseBehaviourFunction, data interface{}) {
		FunctionsMap[FuncName] = Func
		FunctionsData[FuncName] = data
	}
}

/*
func SIM_subprocess_create_async(name string, f BaseBehaviourFunction, host *Host, data interface{}) *Process {
	return SIM_subprocess_create_with_args_async(name, f, host, data, nil)
}

func SIM_subprocess_create_with_args_sync(name string, f BaseBehaviourFunction, host *Host, data interface{}, args []string) *Process {
	p := SIM_subprocess_create_with_args_async(name, f, host, data, args)

	env.stepEnd <- struct{}{}
	<-p.resumeChan
	return p
}

*/
