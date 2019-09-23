package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Printf

func (process *Process) CreateAnomalyHostSync(host *Host, timeEnd float64, newState float64) {
	CreateAnomalyHostAsync(process, host, timeEnd, newState)

	process.env.stepEnd <- struct{}{}
	<-process.resumeChan

	return
}

func CreateAnomalyHostAsync(process *Process, host *Host, timeEnd float64, newState float64) {
	anomaly := NewAnomaly(host, newState)
	event := NewEvent(timeEnd, process, HostAnomalyEvent, nil)
	event.anomaly = anomaly

	heap.Push(&process.env.globalQueue, event)
	return
}

func (process *Process) RepairHostSync(host *Host, timeEnd float64) {
	RepairHostAsync(process, host, timeEnd)

	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func RepairHostAsync(process *Process, host *Host, timeEnd float64) {
	anomaly := NewAnomaly(host, 0)
	event := NewEvent(timeEnd, process, HostRepairEvent, nil)
	event.anomaly = anomaly

	heap.Push(&process.env.globalQueue, event)
}

func breakHost(e *Event) {
	host := e.anomaly.resource.(*Host)
	//heap.Remove(&env.vesninServers, host._index)

	// Break own cores' queue
	coreMap := host.coreManager.coreMap
	for _, core := range coreMap {
		core.breakCore()
	}

	// Break hosts links
	links := GetLinks(host)
	for _, link := range links {
		anom := NewAnomaly(link, 0)
		event := NewEvent(SIM_get_clock(), nil, LinkAnomalyEvent, nil)
		event.anomaly = anom
		breakLink(event)
	}
}

func repairHost(e *Event) {
	host := e.anomaly.resource.(*Host)
	//env.vesninServers.Push(host)

	links := GetLinks(host)

	for _, link := range links {
		link.State = 1.
	}
}

func (core *Core) breakCore() {
	core.state = 0

	for core.wQueue.Len() > 0 {
		item := heap.Pop(&core.wQueue).(*Event)
		item.status = FAIL
		GetCallbacks(item.eventType)[0](item)
	}

	for core.eQueue.Len() > 0 {
		item := heap.Pop(&core.eQueue).(*Event)
		item.status = FAIL
		GetCallbacks(item.eventType)[0](item)
	}
}
