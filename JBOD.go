package lib

import (
	"container/heap"
	"fmt"
	"strconv"
)

var _ = fmt.Print

type JBOD struct {
	Name    string
	Disks   JBODDisks
	DiskArr map[int]*Storage
}

func NewJBOD(JBODId string, amount int, sType string, storageTypes map[string]*StorageType) *JBOD {
	arr := make(JBODDisks, 0)
	heap.Init(&arr)

	jbod := &JBOD{
		Name:  JBODId,
		Disks: arr,
	}

	storageMap := make(map[int]*Storage)

	for j := 0; j < amount; j++ {
		// Name will be like JBOD#1_#4
		storageName := JBODId + "_" + strconv.Itoa(j)
		storage := NewStorage(storageTypes[sType], storageName, jbod)
		storageMap[j] = storage
	}
	jbod.DiskArr = storageMap
	return jbod
}

// Getters and setters
func (jbod *JBOD) GetName() string {
	return jbod.Name
}

func (storage *Storage) GetJBOD() *JBOD {
	return storage.JBOD
}

func getAllJBODsFactory(JBODMap map[string]*JBOD) func() map[string]*JBOD {
	return func() map[string]*JBOD {
		return JBODMap
	}
}

func getCacheJbodFactory(cacheJbod *JBOD) func() *JBOD {
	return func() *JBOD {
		return cacheJbod
	}
}

func getAllJBODsSliceFactory(jbodSlice []*JBOD) func() []*JBOD {
	return func() []*JBOD {
		return jbodSlice
	}
}

/*
=================================================================================================
Storage heap queue (in JBOD) of healthy disks
====================================================================================================
*/

var _ heap.Interface = (*JBODDisks)(nil)

type JBODDisks []*Storage

func (JD JBODDisks) Len() int { return len(JD) }

func (JD JBODDisks) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return JD[i].usedSize < JD[j].usedSize
}

func (JD JBODDisks) Swap(i, j int) {
	JD[i], JD[j] = JD[j], JD[i]
	JD[i]._index = i
	JD[j]._index = j
}

func (JD *JBODDisks) Push(e interface{}) {
	n := len(*JD)
	item := e.(*Storage)
	item._index = n
	*JD = append(*JD, item)
}

func (JD *JBODDisks) Pop() interface{} {
	old := *JD
	n := len(old)
	item := old[n-1]
	item._index = -1 // for safety
	*JD = old[0 : n-1]
	return item
}
