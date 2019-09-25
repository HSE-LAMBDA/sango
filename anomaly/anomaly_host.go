package anomaly

import (
	"container/heap"
	"fmt"
	lib "gosan"
)

var _ = fmt.Printf

func CreateAnomalyHostSync(process *lib.Process, host *lib.Host, timeEnd float64, newState float64) {
	CreateAnomalyHostAsync(process, host, timeEnd, newState)

	process.GetEnv().GetStepEnd() <- struct{}{}
	<-process.GetResumeChan()

	return
}

func CreateAnomalyHostAsync(process *lib.Process, host *lib.Host, timeEnd float64, newState float64) {
	anomaly := NewAnomaly(host, newState)
	event := lib.NewEvent(timeEnd, process, lib.HostAnomalyEvent, nil)
	event.anomaly = anomaly

	heap.Push(process.GetEnv().GetGlobalQueue(), event)
	return
}

func RepairHostSync(process *lib.Process, host *lib.Host, timeEnd float64) {
	RepairHostAsync(process, host, timeEnd)

	process.GetEnv().GetStepEnd()<- struct{}{}
	<-process.GetResumeChan()
}

func RepairHostAsync(process *lib.Process, host *lib.Host, timeEnd float64) {
	anomaly := NewAnomaly(host, 0)
	event := lib.NewEvent(timeEnd, process, lib.HostRepairEvent, nil)
	event.anomaly = anomaly

	heap.Push(process.GetEnv().GetGlobalQueue(), event)
}

//func breakHost(e *lib.Event) {
//	host := e.anomaly.resource.(*lib.Host)
//	//heap.Remove(&env.vesninServers, host._index)
//
//	// Break own cores' queue
//	coreMap := host.coreManager.coreMap
//	for _, core := range coreMap {
//		core.breakCore()
//	}
//
//	// Break hosts links
//	links := lib.GetLinks(host)
//	for _, link := range links {
//		anom := NewAnomaly(link, 0)
//		event := lib.NewEvent(SIM_get_clock(), nil, LinkAnomalyEvent, nil)
//		event.anomaly = anom
//		breakLink(event)
//	}
//}
//
//func repairHost(e *lib.Event) {
//	host := e.anomaly.resource.(*Host)
//	//env.vesninServers.Push(host)
//
//	links := GetLinks(host)
//
//	for _, link := range links {
//		link.State = 1.
//	}
//}
//
//func breakCore(core *lib.Core) {
//	core.state = 0
//
//	for core.wQueue.Len() > 0 {
//		item := heap.Pop(&core.wQueue).(*Event)
//		item.status = FAIL
//		GetCallbacks(item.eventType)[0](item)
//	}
//
//	for core.eQueue.Len() > 0 {
//		item := heap.Pop(&core.eQueue).(*Event)
//		item.status = FAIL
//		lib.GetCallbacks(item.eventType)[0](item)
//	}
//}
