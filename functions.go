package lib

import (
	"container/heap"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var _ = strconv.Atoi
var _ = fmt.Println

func Maximum(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// "component_filename"
func (process *Process) SendPacket(packet *Packet, destinationProcess string) STATUS {
	process._sendPacket(packet, destinationProcess, SyncNetworkEvent)
	return process._add_sync()
}

func (process *Process) DetachedSendPacket(packet *Packet, destinationProcess string) interface{} {
	return process._sendPacket(packet, destinationProcess, AsyncNetworkEvent)
}

func (process *Process) _sendPacket(packet *Packet, destinationProcess string, eventType EventType) interface{} {
	link := process.findLink(destinationProcess)

	timeEnd, _ := getTransferTime(link, packet)

	event := NewEvent(timeEnd, process, eventType, packet)
	event.destinationProcess = destinationProcess

	link.Put(event, &process.env.globalQueue)
	return nil
}

func (process *Process) ReceivePacket(address string) (*Packet, STATUS) {
	defer delete(process.env.workerListeners, address)

	process.env.workerListeners[address] = process.pid

	status := process._add_sync()
	packet := process.env.currentEvent.packet

	return packet, status
}

func (process *Process) SIM_wait(waitTime float64) STATUS {
	core := process.host.coreManager.Next()

	event := NewEvent(process.env.currentTime+waitTime, process, WaitEvent, nil)
	event.resource = core

	heap.Push(&process.env.globalQueue, event)
	heap.Push(&core.wQueue, event)

	return process._add_sync()
}

func (process *Process) findLink(destinationProcess string) *Link {
	i := strings.Index(destinationProcess, "_")

	if i < 0 {
		panic("Index (destination) not found")
	}

	componentName := destinationProcess[:i]
	destinationHost := GetHostByName(componentName)
	link := GetLinkBetweenHosts(process.host, destinationHost)
	return link
}

func (process *Process) SendToHostWithoutReceive(destinationHost *Host, packet *Packet) STATUS {
	link := GetLinkBetweenHosts(process.host, destinationHost)

	timeEnd, _ := getTransferTime(link, packet)

	event := NewEvent(timeEnd, process, SingleNetworkEvent, packet)
	link.Put(event, &process.env.globalQueue)

	return process._add_sync()
}

func (process *Process) _add_sync() STATUS {
	process.env.stepEnd <- struct{}{}
	return <-process.resumeChan
}

func _basic_event_adding_factory() func(*Process, Resource, *Packet, EventType) (STATUS, float64) {
	timeMap := make(map[EventType]func(Resource, *Packet) (float64, float64))
	emptyFunc := func(Resource, *Packet) (float64, float64) { return 0, 0 }

	timeMap[SyncNetworkEvent] = getTransferTime
	timeMap[AsyncNetworkEvent] = getTransferTime
	timeMap[SyncStorageEvent] = getDiskOperationTime
	timeMap[AsyncStorageEvent] = getDiskOperationTime
	timeMap[WaitEvent] = emptyFunc // todo
	timeMap[ExecuteEvent] = getExecutionTime
	timeMap[RecoveryEvent] = getRecoveryTime

	timeMap[LinkAnomalyEvent] = emptyFunc // todo
	timeMap[HostAnomalyEvent] = emptyFunc // todo

	timeMap[LinkRepairEvent] = emptyFunc // todo
	timeMap[HostRepairEvent] = emptyFunc // todo

	timeMap[DiskAnomalyEvent] = emptyFunc // todo

	timeMap[StopEvent] = emptyFunc          // todo
	timeMap[SingleNetworkEvent] = emptyFunc // todo

	return func(process *Process, resource Resource, packet *Packet, eventType EventType) (STATUS, float64) {

		findTimeEnd, ok := timeMap[eventType]
		if !ok {
			log.Panicf("Cannot find time end function for %d\n", eventType)
		}
		timeEnd, t := findTimeEnd(resource, packet)
		event := NewEvent(timeEnd, process, eventType, packet)
		resource.Put(event, &process.env.globalQueue)
		return OK, t
	}
}


type DCAble interface {
	Break(*Process, float64, float64)
	Repair(*Process, float64)
	Update(map[string]float64)
	Reset()

	GetCurrentState() string
	SetCurrentState(string)
}