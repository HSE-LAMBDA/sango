package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Print

type (
	STATUS        int
	EventType     int
	NetworkPacket int
	StoragePacket int
)

const (
	SyncNetworkEvent EventType = iota
	AsyncNetworkEvent
	SyncStorageEvent
	AsyncStorageEvent
	WaitEvent
	ExecuteEvent
	RecoveryEvent

	LinkAnomalyEvent
	HostAnomalyEvent

	LinkRepairEvent
	HostRepairEvent

	DiskAnomalyEvent

	StopEvent
	SingleNetworkEvent
)

const (
	OK STATUS = iota
	FAIL
)

/* ==================================================================

Objects

=====================================================================
*/
func initCallbacksFactory(env *Environment) func(EventType) []func(*Event) {
	callbacks := make(map[EventType][]func(e *Event))

	// Populate Next Worker function

	oneWorkerPopulateFunction := func(e *Event) {
		if e.process != nil {
			e.process.status = e.status
			env.nextWorkers = append(env.nextWorkers, e.process)
		}
	}

	deleteFromResourceLink := func(e *Event) {
		heap.Pop(&e.resource.(*Link).queue)
	}

	deleteFromResourceWait := func(e *Event) {
		heap.Pop(&e.resource.(*Core).wQueue)
	}

	deleteFromResourceExec := func(e *Event) {
		heap.Pop(&e.resource.(*Core).eQueue)
	}

	stopSimulation := func(_ *Event) {
		env.shouldStop = true
	}

	callbacks[SyncNetworkEvent] = append(callbacks[SyncNetworkEvent], func(e *Event) {
		lPID, ok := env.workerListeners[e.destinationProcess]
		if !ok {
			panic(fmt.Sprintf("no such listeners %s", e.destinationProcess))
		}
		receiver, ok := env.workers[lPID]
		if !ok {
			panic(fmt.Sprintf("no such pid"))
		}
		receiver.status = e.status
		e.process.status = e.status

		env.nextWorkers = append(env.nextWorkers, e.process, receiver)
	},
		deleteFromResourceLink)

	callbacks[AsyncNetworkEvent] = append(callbacks[AsyncNetworkEvent], func(e *Event) {
		lPID, ok1 := env.workerListeners[e.destinationProcess]
		if !ok1 {
			panic(fmt.Sprintf("no such listener %s", e.destinationProcess))
		}
		receiver, ok2 := env.workers[lPID]
		if !ok2 {
			panic("no such pid")
		}
		receiver.status = e.status
		env.nextWorkers = append(env.nextWorkers, receiver)
	},
		deleteFromResourceLink)

	callbacks[SingleNetworkEvent] = append(callbacks[SingleNetworkEvent], oneWorkerPopulateFunction, deleteFromResourceLink)
	callbacks[SyncStorageEvent] = append(callbacks[SyncStorageEvent], oneWorkerPopulateFunction, deleteFromResourceLink)
	callbacks[AsyncStorageEvent] = append(callbacks[AsyncStorageEvent], func(e *Event) {}, deleteFromResourceLink)
	callbacks[WaitEvent] = append(callbacks[WaitEvent], oneWorkerPopulateFunction, deleteFromResourceWait)
	callbacks[ExecuteEvent] = append(callbacks[ExecuteEvent], oneWorkerPopulateFunction, deleteFromResourceExec)
	callbacks[RecoveryEvent] = append(callbacks[RecoveryEvent], oneWorkerPopulateFunction, deleteFromResourceExec)

	callbacks[StopEvent] = append(callbacks[StopEvent], stopSimulation)

	// Anomaly Event callbacks //todo I deleted oneWorkerPopulateFunction
	callbacks[LinkAnomalyEvent] = append(callbacks[LinkAnomalyEvent], breakLink)
	callbacks[HostAnomalyEvent] = append(callbacks[HostAnomalyEvent], breakHost)

	// Repair Event callbacks
	callbacks[LinkRepairEvent] = append(callbacks[LinkRepairEvent], repairLink)
	callbacks[HostRepairEvent] = append(callbacks[HostRepairEvent], repairHost)

	// Disk Anomalies
	callbacks[DiskAnomalyEvent] = append(callbacks[DiskAnomalyEvent], breakDisk)

	return func(t EventType) []func(*Event) {
		return callbacks[t]
	}
}

type Event struct {
	id        string
	eventType EventType
	timeStart float64
	timeEnd   float64
	status    STATUS
	process   *Process

	index    int
	index_g  int
	priority int

	packet   *Packet
	resource Resource

	destinationProcess string

	anomaly *Anomaly
}

//(link, task, process, destinationProcess)
func NewEvent(timeEnd float64, process *Process, eventType EventType, packet *Packet) *Event {
	e := &Event{
		timeStart: SIM_get_clock(),
		timeEnd:   timeEnd,
		process:   process,
		eventType: eventType,
		packet:    packet,
	}

	return e
}
