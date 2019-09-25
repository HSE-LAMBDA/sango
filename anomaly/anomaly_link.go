package anomaly

import (
	"container/heap"
	"fmt"
	lib "gosan"
)

var _ = fmt.Print

func CreateAnomalyLinkSync(process *lib.Process, linkId *lib.Link, timeClock float64, newState float64) {
	CreateAnomalyLinkAsync(process, linkId, timeClock, newState)
	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func CreateAnomalyLinkAsync(process *lib.Process, link *lib.Link, timeClock float64, newState float64) {
	anom := NewAnomaly(link, newState)
	event := lib.NewEvent(timeClock, process, lib.LinkAnomalyEvent, nil)
	event.anomaly = anom

	heap.Push(&process.env.globalQueue, event)
}

func RepairLinkSync(process *lib.Process, link *lib.Link, timeClock float64, newState float64) {
	RepairLinkAsync(process, link, timeClock)
	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func RepairLinkAsync(process *lib.Process, link *lib.Link, timeClock float64) {
	anom := NewAnomaly(link, 1)
	event := lib.NewEvent(timeClock, process, lib.LinkRepairEvent, nil)
	event.anomaly = anom

	x := process.env.globalQueue
	_ = x
	heap.Push(&process.env.globalQueue, event)
}

func breakLink(e *lib.Event) {
	link := e.anomaly.resource.(*lib.Link)
	newState := e.anomaly.newState

	oldState := link.State
	link.State = newState

	if newState == 0 {
		for len(link.queue) > 0 {
			event := heap.Pop(&link.queue).(*lib.Event)
			event.status = FAIL
			lib.GetCallbacks(event.eventType)[0](event)
		}
	} else {
		for _, event := range link.queue {
			event.timeEnd = oldState * (event.timeEnd - event.timeStart) / link.State
		}
		heap.Init(&e.process.env.globalQueue)
	}

}

func repairLink(e *lib.Event) {
	link := e.anomaly.resource.(*lib.Link)
	link.State = 1.
}
