package src

import (
	"lib"
	"log"
	"strconv"
)

type Client struct {
	file      *lib.File
	timeStart float64
	priority  float64

	index int
}

type ClientQueue []*Client

func (cq ClientQueue) Len() int { return len(cq) }

func (cq ClientQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return cq[i].priority > cq[j].priority
}

func (cq ClientQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
	cq[i].index = i
	cq[j].index = j
}

func (cq *ClientQueue) Push(x interface{}) {
	n := len(*cq)
	client := x.(*Client)
	client.index = n
	*cq = append(*cq, client)
}

func (cq *ClientQueue) Pop() interface{} {
	old := *cq
	n := len(old)
	client := old[n-1]
	client.index = -1 // for safety
	*cq = old[0 : n-1]
	return client
}

func (CL *ClientManager) SyncClientManagerProcess(TP *SANProcess, data interface{}) {
	arg, ok := data.(*ClientFlags)
	if !ok {
		log.Panic("Async client flag data conversion error")
	}

	host := CL.Host
	for i := 0; i < arg.FileAmount; i++ {

		fileSize := fRand(arg.MinFileSize, arg.MaxFileSize)
		waitTime := fRand(arg.MinPauseTime, arg.MaxPauseTime)

		file := lib.NewFile(strconv.Itoa(i), fileSize, arg.RequestType, 1, lib.GetRandomBlockSize(), 0)
		CL.FORK_SYNC_CLIENT("Client_"+file.Filename, CL.PacketSenderReceiverProcess, host, file)
		TP.SIM_wait(waitTime)
	}
	return
}
func (CL *ClientManager) FORK_SYNC_CLIENT(name string, f SANBFunction, host *lib.Host, file *lib.File) {
	switch file.RequestType {
	case lib.WRITE, lib.RANDWRITE:
		CL.iob.WriteQueueLength += 1
	case lib.RANDREAD, lib.READ:
		CL.iob.ReadQueueLength += 1
	}

	FORK_SYNC("Client_"+file.Filename, CL.PacketSenderReceiverProcess, host, file)

	switch file.RequestType {
	case lib.RANDWRITE, lib.WRITE:
		CL.iob.WriteQueueLength -= 1
	case lib.RANDREAD, lib.READ:
		CL.iob.ReadQueueLength -= 1
	}
}
