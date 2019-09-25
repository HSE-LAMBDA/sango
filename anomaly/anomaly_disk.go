package anomaly

import (
	"container/heap"
	"fmt"
	lib "gosan"
)

var _ = fmt.Print

// JBOD#1_#4

//func (process *lib.Process) CreateAnomalyDiskSync(storage *lib.Storage, timeClock float64, anomalyPart float64) {
func CreateAnomalyDiskSync(process *lib.Process, storage *lib.Storage, timeClock float64, anomalyPart float64) {
	CreateDiskAnomalyAsync(process, storage, timeClock, anomalyPart)
	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func CreateDiskAnomalyAsync(process *lib.Process, storage *lib.Storage, timeClock float64, anomalyPart float64) {
	anom := NewAnomaly(storage, anomalyPart)

	event := lib.NewEvent(timeClock, process, lib.DiskAnomalyEvent, nil)
	event.anomaly = anom

	heap.Push(&process.env.globalQueue, event)

}

func CreateDiskAnomalyFullOFF(process *lib.Process, timeClock float64, storage *lib.Storage, anomalyPart float64) {
	// Полный выход диска из строя

	anom := &Anomaly{
		resource:    storage,
		anomalyPart: anomalyPart,
	}

	event := *lib.NewEvent(timeClock, process, DiskAnomalyEvent, nil)
	event.anomaly = anom

	heap.Push(&process.GetEnv().globalQueue, event)

	process.env.stepEnd <- struct{}{}
	<-process.resumeChan
}

func breakDisk(e *lib.Event) {
	disk := e.anomaly.resource.(*Storage)
	anomalyPart := e.anomaly.anomalyPart

	disk.brokenPart = anomalyPart

	if anomalyPart == 1 {
		disk.fail = true

		wQueue := disk.WriteLink.queue
		for len(wQueue) > 0 {
			event := heap.Pop(&wQueue).(*Event)
			event.status = FAIL
			GetCallbacks(event.eventType)[0](event)
		}

		rQueue := disk.ReadLink.queue
		for len(rQueue) > 0 {
			event := heap.Pop(&rQueue).(*Event)
			event.status = FAIL
			GetCallbacks(event.eventType)[0](event)
		}

	} else {
		// todo
	}

}

func repairDisk(e *lib.Event) {
	link := e.anomaly.resource.(*lib.Link)
	link.State = 1.
}
