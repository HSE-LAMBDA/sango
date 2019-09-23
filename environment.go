package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Print

type Environment struct {
	globalQueue     globalEventQueue
	currentEvent    *Event
	currentTime     float64
	workers         map[uint64]*Process
	workerListeners map[string]uint64
	daemonList      map[uint64]*Process
	shouldStop      bool
	stepEnd         chan interface{}
	nextWorkers     []*Process
	systemResources []Closeable
}

func NewEnvironment() *Environment {
	queue := make(globalEventQueue, 0)
	heap.Init(&queue)

	e := &Environment{
		globalQueue:     queue,
		workers:         make(map[uint64]*Process),
		daemonList:      make(map[uint64]*Process),
		workerListeners: make(map[string]uint64),
		nextWorkers:     make([]*Process, 0),
		stepEnd:         make(chan interface{}),
	}

	return e
}

func (env *Environment) stopSimulation() {
	env.shouldStop = true
}

func (env *Environment) updateQueue(deltaTime float64) {

}

func (env *Environment) SendResumeSignalToWorkers() {
	for len(env.nextWorkers) > 0 {
		env.nextWorkers[0].resumeChan <- env.nextWorkers[0].status
		<-env.stepEnd
		env.nextWorkers = env.nextWorkers[1:]
	}
}

/*
 1) Собрать минимальные события со всех компонент
 2) Выбрать текущее событие
 3) Отправить сигнал компоненте, что её событие минимально
 4) Компонента должна понять, какой LP породил это событие и отправить ему сигнал
*/

func (env *Environment) Step() {

	wLen := len(env.workers)
	qLen := len(env.globalQueue)
	dLen := len(env.daemonList)

	//if (wLen-dLen == 0 && qLen == 0) || (wLen == dLen) && qLen == wLen {
	if wLen == 0 && qLen == 0 || qLen == dLen {
		if qLen > 0 {
			env.currentTime = env.globalQueue[0].timeEnd
		}
		env.shouldStop = true
		return
	}

	if qLen == 0 {
		panic("deadlock")
	}

	event := heap.Pop(&env.globalQueue).(*Event)

	if event.status == FAIL {
		return
	}
	env.update(event)

	// execute callbacks 1) delete self from source and 2) populate next workers slice
	callbacks := GetCallbacks(event.eventType)
	for _, function := range callbacks {
		function(event)
	}

}

func (env *Environment) update(currentEvent *Event) {
	env.updateQueue(currentEvent.timeEnd - env.currentTime)
	env.currentTime = currentEvent.timeEnd
	env.currentEvent = currentEvent
}
