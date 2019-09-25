package lib

import (
	"fmt"
)

var _ = fmt.Fprintf

type (
	ProcessInterface interface {
		SetParent(*Process)
		GetParent() *Process
	}
	MachineInterface interface {
	}

	Process struct {
		env        *Environment
		resumeChan chan STATUS
		host       *Host

		name string
		data interface{}

		pid    uint64
		status STATUS
	}

	ProcessID struct {
		id uint64
	}

	Closeable interface {
		Close()
	}
)

func PIDNextFactory() func() uint64 {
	var i uint64
	return func() uint64 {
		i++
		return i
	}
}

func NewProcessFactory(env *Environment) func(string, *Host) *Process {
	return func(name string, host *Host) *Process {
		pid := PIDNext()
		p := &Process{
			name:       name,
			env:        env,
			resumeChan: make(chan STATUS),
			host:       host,
			pid:        pid,

			status: OK,
		}

		env.workers[pid] = p
		env.nextWorkers = append(env.nextWorkers, p)
		return p
	}
}

func ProcWrapper(processStrategy BaseBehaviourFunction, pi ProcessInterface, data interface{}) {
	go func() {
		p := pi.GetParent()
		<-p.resumeChan
		processStrategy(pi, data)
		delete(p.env.workers, p.pid)
		p.env.stepEnd <- struct{}{}
	}()
}

func ProcWrapperTemp(processStrategy func(*Process, interface{}), p *Process, data interface{}) {
	go func() {
		<-p.resumeChan
		processStrategy(p, data)
		delete(p.env.workers, p.pid)
		p.env.stepEnd <- struct{}{}
	}()
}

func (process *Process) Daemonize(resource ...Closeable) {
	process.env.daemonList[process.pid] = process
	process.env.systemResources = append(process.env.systemResources, resource...)
}

func (process *Process) GetPID() uint64 {
	return process.pid
}

func (process *Process) SetHost(host *Host) {
	process.host = host
}

func (process *Process) Execute(packet *Packet) (STATUS, float64) {
	core := process.host.coreManager.Next()
	_, t := _basic_event_adding(process, core, packet, ExecuteEvent)
	return process._add_sync(), t
}

func (process *Process) RecoveryExecutePacket(packet *Packet) (*Packet, STATUS) {
	core := process.host.coreManager.Next()
	_basic_event_adding(process, core, packet, RecoveryEvent)

	return PACKET_RECOVERY, process._add_sync()
}

func (process *Process) GetData() interface{} {
	return process.data
}

func (process *Process) GetName() string {
	return process.name
}

func (process *Process) GetEnv() *Environment {
	return process.env
}


func (process *Process) GetResumeChan() chan STATUS{
	return process.resumeChan
}