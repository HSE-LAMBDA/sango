package lib

import (
	"container/heap"
	"fmt"
	"log"
)

var _ = fmt.Print
var _ = heap.Init

type Link struct {
	cid       uint64
	Name      string
	Bandwidth float64
	Latency   float64
	State     float64
	queue     EventQueue

	lastTime float64

	Src, Dst *Host
}

func NewLink(bandwidth float64, name string) *Link {
	queue := make(EventQueue, 0)
	heap.Init(&queue)
	link := &Link{
		cid:       CIDNext(),
		queue:     queue,
		Bandwidth: bandwidth,
		Name:      name,
		State:     1.,
	}
	return link
}

func (link *Link) Put(e *Event, globalQueue *globalEventQueue) {
	e.resource = link
	if link.State == 0 {
		e.timeEnd = SIM_get_clock()

		e.status = FAIL
		GetCallbacks(e.eventType)[0](e)
		return
	}
	link.lastTime = e.timeEnd

	heap.Push(globalQueue, e)
	heap.Push(&link.queue, e)
}

func getLinkBetweenHostsFactory(routesMap, backupRoutesMap map[*Host]map[*Host]*Link) func(*Host, *Host) *Link {
	return func(start *Host, finish *Host) *Link {
		hostLinkMap, ok := routesMap[start]
		if !ok {
			log.Panic("No such host ", start.Name)
		}
		link, ok := hostLinkMap[finish]
		if ok {
			return link
		}

		// look for backup routes
		hostLinkMap, ok = backupRoutesMap[start]
		if !ok {
			log.Panic("No such host ", start.Name)
		}
		link, ok = hostLinkMap[finish]
		if ok {
			return link
		}
		log.Panic(fmt.Sprintf("No such route between hosts: %s %s", start.Name, finish.Name))
		return nil
	}
}

func getLinksFactory(hostsLinksMap map[*Host][]*Link) func(*Host) []*Link {
	return func(host *Host) []*Link {
		return hostsLinksMap[host]
	}
}

func (link *Link) getCID() uint64 {
	return link.cid
}

func (link *Link) GetQueueSize() int {
	return len(link.queue)
}

func getLinkByNameFactory(linksMap map[string]*Link) func(string) *Link {
	return func(name string) *Link {
		return linksMap[name]
	}
}

func getAllLinksMapFactory(linksMap map[string]*Link) func() map[string]*Link {
	return func() map[string]*Link {
		return linksMap
	}
}
