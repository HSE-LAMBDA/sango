package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Println

type (
	CoreManager struct {
		// Round Robin implementation
		id    uint64
		start uint64
		end   uint64

		coreMap map[uint64]*Core
	}

	Host struct {
		Name        string
		Type        string
		Id          string
		Speed       float64
		processes   map[uint64]*Process
		storage     map[int]*Storage
		cid         uint64
		nCore       uint64
		coreManager *CoreManager
	}

	Core struct {
		id       uint64
		state    float64
		speed    float64
		eQueue   EventQueue
		wQueue   EventQueue
		lastTime float64
	}
)

func NewCoreManager(start uint64, end uint64, speed float64, state float64) *CoreManager {
	m := make(map[uint64]*Core)

	for i := start; i < end; i++ {
		id := CIDNext()
		eQueue := make(EventQueue, 0)
		wQueue := make(EventQueue, 0)

		heap.Init(&eQueue)
		heap.Init(&wQueue)

		core := &Core{
			id:     id,
			eQueue: eQueue,
			wQueue: wQueue,
			state:  state,
			speed:  speed,
		}

		m[i] = core
	}
	c := &CoreManager{
		start: start,
		id:    start,
		end:   end,

		coreMap: m,
	}
	return c
}

func (cid *CoreManager) Next() *Core {

	if cid.id+1 < cid.end {
		cid.id++
	} else {
		cid.id = cid.start
	}
	return cid.coreMap[cid.id]
}

func NewHost(name, typeID, ID string, speed float64, nCore uint64) *Host {

	coreManager := NewCoreManager(0, nCore, speed, 1)

	host := &Host{
		cid:   CIDNext(),
		Name:  name,
		Id:    ID,
		Type:  typeID,
		Speed: speed,
		nCore: nCore,

		coreManager: coreManager,
	}
	return host
}

func (host *Host) Put(e *Event, globalQueue *globalEventQueue) {

}

func (core *Core) Put(e *Event, globalQueue *globalEventQueue) {
	e.resource = core
	core.lastTime = e.timeEnd

	heap.Push(globalQueue, e)
	heap.Push(&core.eQueue, e)
}

func (host *Host) getCID() uint64 {
	return host.cid
}

func (core *Core) getCID() uint64 {
	return core.id
}
