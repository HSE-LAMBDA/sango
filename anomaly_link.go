package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Print

func (process *Process) CreateAnomalyLinkSync(linkId *Link, timeClock float64, newState float64) {
	CreateAnomalyLinkAsync(process, linkId, timeClock, newState)
	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func CreateAnomalyLinkAsync(process *Process, link *Link, timeClock float64, newState float64) {
	anom := NewAnomaly(link, newState)
	event := NewEvent(timeClock, process, LinkAnomalyEvent, nil)
	event.anomaly = anom

	heap.Push(&process.env.globalQueue, event)
}

func (process *Process) RepairLinkSync(link *Link, timeClock float64, newState float64) {
	RepairLinkAsync(process, link, timeClock)
	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func RepairLinkAsync(process *Process, link *Link, timeClock float64) {
	anom := NewAnomaly(link, 1)
	event := NewEvent(timeClock, process, LinkRepairEvent, nil)
	event.anomaly = anom

	x := process.env.globalQueue
	_ = x
	heap.Push(&process.env.globalQueue, event)
}

func breakLink(e *Event) {
	link := e.anomaly.resource.(*Link)
	newState := e.anomaly.newState

	oldState := link.State
	link.State = newState

	if newState == 0 {
		for len(link.queue) > 0 {
			event := heap.Pop(&link.queue).(*Event)
			event.status = FAIL
			GetCallbacks(event.eventType)[0](event)
		}
	} else {
		for _, event := range link.queue {
			event.timeEnd = oldState * (event.timeEnd - event.timeStart) / link.State
		}
		heap.Init(&e.process.env.globalQueue)
	}

}

func repairLink(e *Event) {
	link := e.anomaly.resource.(*Link)
	link.State = 1.
}
