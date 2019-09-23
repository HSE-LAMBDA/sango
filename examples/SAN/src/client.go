package src

import (
	"lib"
	"log"
	"math"
	"math/rand"
	"strconv"
)

type (
	ClientManager struct {
		*lib.Host  `json:"-"`
		iob        *IOBalancer     `json:"-"`
		Cap        float64         `json:"cap"`
		frequency  float64         `json:"-"`
		id         string          `json:"-"`
		numJobsMap map[int]float64 `json:"-"`
	}
	Master struct {
		*lib.Host
		clients map[int]*ClientManager
	}
)

func NewClientManager(host *lib.Host, iob *IOBalancer, id string, capacity float64,
	parMap map[int]float64) *ClientManager {

	return &ClientManager{
		Host: host,
		iob:  iob,
		Cap:  capacity,
		id:   id,
		numJobsMap:parMap,
	}
}

func NewMaster(host *lib.Host, clients map[int]*ClientManager) *Master {
	return &Master{
		Host:    host,
		clients: clients,
	}
}

func (M *Master) ClientManagerAsyncProcess(TP *SANProcess, data interface{}) {
	arg, ok := data.(*ClientFlags)
	if !ok {
		log.Panic("Async client flag data conversion error")
	}

	host := M.Host
	fileSize := fRand(arg.MinFileSize, arg.MaxFileSize)
	waitTime := fRand(arg.MinPauseTime, arg.MaxPauseTime)

	for i := 0; i < arg.FileAmount; i++ {

		numJobs := len(M.clients)
		fileSize /= float64(numJobs)

		for _, client := range M.clients {
			if numJobs <= 0 {
				break
			}
			index := strconv.Itoa(i) + "_" + strconv.Itoa(numJobs)
			file := lib.NewFile(index, fileSize, arg.RequestType, len(M.clients), lib.GetRandomBlockSize(), 0)
			FORK("Client_"+index, client.PacketSenderReceiverProcess, host, file)
			numJobs--
		}
		TP.SIM_wait(waitTime)
	}
	return
}

func (CL *ClientManager) PacketSenderReceiverProcess(TP *SANProcess, data interface{}) {
	iob := CL.iob
	CL.iob.IOProcessCnt++
	defer CL.DecreaseIOCNT()

	file, ok := data.(*lib.File)
	if !ok {
		log.Panic("Client file conversion error")
	}

	TP.SIM_wait(file.Latency)

	fileInfo := iob.CreateSANExecutors(file)

	packetSize := fileInfo.Packet1.Size
	cap := CL.Cap

	nSending := file.Size / file.Packet1.Size
	//nSending = 1

	for i := 0.; i < nSending; i++ {
		res := CL.SendPacket(TP, fileInfo)
		if res != lib.OK {
			return
		}

		res = CL.ReceivePacket(TP, fileInfo)
		if res != lib.OK {
			return
		}

		// Like a queue
		parallelOverhead, ok := CL.numJobsMap[file.NumJobs]
		if ok == false {
			panic("no such numjobs overhead time")
		}

		//parallelOverhead = math.Pow(parallelOverhead, float64(file.NumJobs)) / float64(file.NumJobs)
		parallelOverhead = parallelOverhead * math.Exp(0.1 * float64(file.NumJobs)) / float64(file.NumJobs)
		TP.SIM_wait(parallelOverhead + (packetSize / cap))

	}

	log.Printf("File%s,%.3f, %.2f, %s\n", file.Filename, file.Size/1e6, lib.SIM_get_clock(), fileInfo.serverName)
	TP.SendPacket(lib.PACKET_FINALIZE, fileInfo.clientContr)
}

func (CL *ClientManager) WriteProcessRestart(fileInfo *FileInfo) {
	CL.iob.GlobalDiskState.Crrq++
	switch fileInfo.RequestType {
	case lib.WRITE, lib.RANDWRITE:
		CL.iob.WriteCancelRate++
	case lib.RANDREAD, lib.READ:
		CL.iob.ReadCancelRate++
	}

	FORK("", CL.PacketSenderReceiverProcess, CL.Host, fileInfo)
	//fmt.Println("Restart", fileInfo.serverName, lib.SIM_get_clock(), fileInfo.Filename)
}

func fRand(fMin float64, fMax float64) float64 {
	return fMin + rand.Float64()*(fMax-fMin)
}

func (CL *ClientManager) DecreaseIOCNT() {
	CL.iob.IOProcessCnt--
}

func (CL *ClientManager) Tracer(TP *SANProcess, _ interface{}) {
	return
	TP.Daemonize()
	for {
		CL.IOBalancerWriteCounterProcess(TP, nil)
		CL.IOBalancerReadCounterProcess(TP, nil)
		TP.SIM_wait(0.1)
	}
	return
}

func (CL *ClientManager) SendPacket(TP *SANProcess, info *FileInfo) lib.STATUS {
	info.metric.timeStart = lib.SIM_get_clock()
	CL.iob.packetCounter++
	res := TP.SendPacket(info.Packet1, info.clientContr)

	if res != lib.OK {
		CL.WriteProcessRestart(info)
		return res
	}

	return res
}

func (CL *ClientManager) ReceivePacket(TP *SANProcess, info *FileInfo) lib.STATUS {
	_, res := TP.ReceivePacket(info.contrClient)
	if res != lib.OK {
		CL.WriteProcessRestart(info)
	}

	iob := CL.iob

	t2 := lib.SIM_get_clock()
	deltaTime := t2 - info.metric.timeStart

	switch info.RequestType {
	case lib.WRITE, lib.RANDWRITE:
		t := iob.WriteResponseRate + (deltaTime-iob.WriteResponseRate)/iob.packetCounter
		iob.WriteResponseRate = t
		iob.GlobalDiskState.WriteAwaitMs = t
	case lib.READ, lib.RANDREAD:
		t := iob.ReadResponseRate + (deltaTime-iob.ReadResponseRate)/iob.packetCounter
		iob.ReadResponseRate = t
		iob.GlobalDiskState.ReadAwaitMs = t

		//fmt.Println(t)
	}

	return res
}
