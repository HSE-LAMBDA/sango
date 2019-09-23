package lib

import (
	"container/heap"
	"fmt"
)

var _ = fmt.Print

/*
=================================================================================================
EventInterface event queue
====================================================================================================
*/

var _ heap.Interface = (*EventQueue)(nil)

type EventQueue []*Event

func (eq EventQueue) Len() int { return len(eq) }

func (eq EventQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return eq[i].timeEnd < eq[j].timeEnd
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].index = i
	eq[j].index = j
}

func (eq *EventQueue) Push(e interface{}) {
	n := len(*eq)
	item := e.(*Event)
	item.index = n
	*eq = append(*eq, item)
}

func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*eq = old[0 : n-1]
	return item
}

/*
=================================================================================================
Global event queue
====================================================================================================
*/

var _ heap.Interface = (*globalEventQueue)(nil)

type globalEventQueue []*Event

func (eq globalEventQueue) Len() int { return len(eq) }

func (eq globalEventQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return eq[i].timeEnd < eq[j].timeEnd
}

func (eq globalEventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].index_g = i
	eq[j].index_g = j
}

func (eq *globalEventQueue) Push(e interface{}) {
	n := len(*eq)
	item := e.(*Event)
	item.index_g = n
	*eq = append(*eq, item)
}

func (eq *globalEventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	item := old[n-1]
	item.index_g = -1 // for safety
	*eq = old[0 : n-1]
	return item
}
