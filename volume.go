package lib

import "encoding/xml"

type Volume struct {
	Id     string       `xml:"id,attr"`
	Mounts []VolumePart `xml:"mount"`
}

type VolumePart struct {
	XMLName xml.Name `xml:"mount"`
	Id      string   `xml:"id,attr"`
	// Fraction float64  `xml:"fraction,attr"`
}

func getAllVolumesFactory(volumes []*Volume) func() map[string]*Volume {
	vMap := make(map[string]*Volume)
	for _, vol := range volumes {
		vMap[vol.Id] = vol
	}
	return func() map[string]*Volume {
		return vMap
	}
}

/*



import (
	"math/rand"
	"fmt"
	"sync/atomic"
)

var _ = fmt.Println

type Volume struct {
	name     string
	storages map[string]*Storage
	index []string
}



func NewVolume(name string, storages map[string]*Storage) (volume *Volume) {
	diskNames := make([]string, len(storages))

	i := 0
	for k := range storages {
		diskNames[i] = k
		i++
	}

	volume = &Volume{
		name:     name,
		storages: storages,
		index:diskNames,
	}

	return
}


func (process *Process) WriteToVolume(volume *Volume, data *Task) interface{} {
	/*
		- Create a list of storage disks write to
		- Call a library function and write

		hard code:
		8+2 parities todo


	numberOfFragments := 5

	storageAmount := len(volume.storages)

	if storageAmount <= numberOfFragments{
		numberOfFragments = storageAmount
	}

	diskIndexes := sampleConsequently(storageAmount, numberOfFragments)

	currentLogicBlockAmount := uint8(0)

	storageBlockTraceHelper := &StorageBlockTraceHelper{
		currentBlockAmount:&currentLogicBlockAmount,
		totalBlockAmount:uint8(numberOfFragments),
	}

	task := NewTask("logicwrite", 0, data.size/float64(numberOfFragments), storageBlockTraceHelper)

	// Write
	for i := range diskIndexes{
		process.XWrite(volume.storages[volume.index[i]], task)
	}

	// End of my operation
	env.stepEnd <- struct{}{}

	// Wait until all
	<-process.resumeChan
	return nil
}

// Mix of synchronous and async write operation
func (process *Process) XWrite(storage *Storage, task *Task) interface{} {
	storageBlockTraceHelper := task.data.(*StorageBlockTraceHelper)

	event := NewSendEvent(worker, task, "")
	event.storage = true
	event.currentLogicBlockAmount = storageBlockTraceHelper.currentBlockAmount
	event.totalBlockAmount = storageBlockTraceHelper.totalBlockAmount
	event.link = storage.writeLink


	atomic.AddInt64(&storage.writeLink._counter, 1)

	t := &TransferEvent{
		sendEvent:    event,
		receiveEvent: nil,
	}
	event.callbacks = append(event.callbacks, deleteSelfFromResource)
	storage.writeLink.Put(t)
	return nil
}


// Mix of synchronous and async reading
func (process *Process) XRead(storage *Storage, task *Task) interface{} {
	event := NewSendEvent(worker, task, "")
	event.storage = true
	event.link = storage.readLink
	atomic.AddInt64(&storage.readLink._counter, 1)
	t := &TransferEvent{
		sendEvent:    event,
		receiveEvent: nil}

	event.callbacks = append(event.callbacks, deleteSelfFromResource)
	storage.readLink.Put(t)
	return nil
}


// Return a slice starting at $low with length $n
func sampleConsequently(diskAmount int, n int) []int{
	maxLow := diskAmount - n + 1
	arr := make([]int, n)

 	if maxLow <= 0{
		index := 0
		for ; index < diskAmount; index++{
			arr[index] = index
		}
		for ;index < n; index++{
			arr[index] = rand.Intn(diskAmount)
		}
	}else {
		low := rand.Intn(maxLow)
		for i := 0; i < n; i++{
			arr[i] = low
			low++
		}
	}
	return arr
}


// Return a slice of random int [0, high) with length n
func sampleInt(high int, n int) []int {
	generated := make(map[int]bool)
	if high <= 0 {
		fmt.Println("protpor")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		number := rand.Intn(high)
		if !generated[number]{
			generated[number] = true
			arr[i] = number
		}
	}
	return arr
}


type StorageBlockTraceHelper struct {
	totalBlockAmount uint8
	currentBlockAmount *uint8
}

*/

// TODO Check again. May be incorrect
//if e, ok := currentEvent.(*TTTEvent); ok {
//	if e.storage {
//		*e.currentLogicBlockAmount++
//		if *e.currentLogicBlockAmount < e.totalBlockAmount {
//			goto start
//		}
//	}
//}
