package src

import (
	lib "gosan"
	"log"
)

type (
	SANJBODController struct {
		*lib.Host     `json:"-"`
		*NamingProps  `json:"-"`
		*CommonProps  `json:"-"`
		*SANComponent `json:"-"`
		disks         map[int]*SANDisk `json:"-"`
		disksSlice    []*SANDisk       `json:"-"`
		iob           *IOBalancer      `json:"-"`
	}

	DiskState struct {
		Pool       string `json:"pool"`
		TotalSpace uint64 `json:"total_cap_MB"` // TODO
		UsedSpace  uint64 `json:"used_cap_MB"`  // TODO
		AllocMB    uint64 `json:"alloc_MB"`     // TODO
		LenMB      uint64 `json:"len_MB"`       // TODO
		Status     string `json:"status"`       // TODO
		State      string `json:"state"`        // TODO

		AvgReadSpeed  float64     `json:"r_KBps"` // TODO `json:"avg_read_speed"`
		AvgWriteSpeed float64     `json:"w_KBps"` // TODO `json:"avg_write_speed"`
		ReadOps       float64     `json:"r_ops"`
		WriteOps      float64     `json:"w_ops"`
		ReadAwaitMs   float64     `json:"r_await_ms"`
		WriteAwaitMs  float64     `json:"w_await_ms"`
		Crrm          float64     `json:"crrm"`     // TODO
		Rrqm          float64     `json:"rrqm"`     // TODO
		Wrqm          float64     `json:"wrqm"`     // TODO
		Avgrq_sz      float64     `json:"avgrq_sz"` // TODO NZ
		Avgqu_sz      float64     `json:"avgqu_sz"` // TODO NZ
		MBRead        float64     `json:"MB_read"`
		MBWrtn        float64     `json:"MB_wrtn"`
		IORead        [10]float64 `json:"blk_read"`
		IOWrite       [10]float64 `json:"blk_write"`

		// unnecessary
		Async   float64 `json:"async"`
		Crrq    float64 `json:"crrq"`
		AvgQsz  float64 `json:"avg_qsz"`
		SvctmMs float64 `json:"svctm_ms"`
	}
)

func (jbod *SANJBODController) JbodWriterProcess(TP *SANProcess, data interface{}) {
	arg, ok := data.(*FileInfo)
	if !ok {
		log.Panic("Fileinfo conversion error")
	}

	receiveAddress1 := arg.cacheJbod
	sendAddress1 := arg.jbodCache

	disks := jbod.disksSlice

	diskAmount := len(disks)
	roundRobinIndex := 0

	for {
		task, res := TP.ReceivePacket(receiveAddress1)
		if res != lib.OK || task == lib.PACKET_FINALIZE {
			break
		}

		roundRobinIndex = (roundRobinIndex + 1) % diskAmount
		disks[roundRobinIndex].WriteAsync(TP, task)
	}
	_ = sendAddress1
}

func (jbod *SANJBODController) JbodReaderProcess(TP *SANProcess, data interface{}) {
	arg, ok := data.(*FileInfo)
	if !ok {
		log.Panic("Fileinfo conversion error")
	}
	dataBlock := arg.Packet2

	mailbox := arg.cacheJbod
	sendAddr := arg.jbodCache

	for {
		task, res := TP.ReceivePacket(mailbox)
		if task == lib.PACKET_FINALIZE {
			break
		}
		if res != lib.OK || task == lib.PACKET_FINALIZE {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddr)
			break
		}

		// Найти диск, где лежит пакет с данными
		disk := jbod.GetPacketLocation(task)

		// причесать то, как надо будет доставать dataBlock
		_, res = disk.ReadSync(TP, dataBlock)
		if res != lib.OK {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddr)
			break
		}
		res = TP.SendPacket(dataBlock, sendAddr)
		if res != lib.OK {
			break
		}
	}
}

func (jbod *SANJBODController) GetPacketLocation(packet *lib.Packet) *SANDisk {
	randomNumber := randomInt(0, len(jbod.disksSlice))
	return jbod.disksSlice[randomNumber]
}

// JBOD name matches with the name of its controller
func NewSANJBODController(libJBOD *lib.JBOD) *SANJBODController {
	libDisks := libJBOD.Disks
	SANDisksMap := make(map[int]*SANDisk)
	SANDisksSlice := make([]*SANDisk, len(libDisks))

	host := lib.GetHostByName(libJBOD.Name)

	naming := NewNamingProps(libJBOD.Name, "Storage", libJBOD.Name)
	jbod := &SANJBODController{
		Host:        host,
		NamingProps: naming,
		SANComponent: &SANComponent{
			currentState: "default",
		},
		CommonProps: &CommonProps{
			Status: OK,
		},
	}

	k := 0
	for i, disk := range libDisks {
		tDisk := NewSANDisk(disk, jbod)
		tDisk.index = k
		SANDisksMap[i] = tDisk
		SANDisksSlice[k] = tDisk
		k++
	}

	jbod.disks = SANDisksMap
	jbod.disksSlice = SANDisksSlice

	return jbod
}

func (tjc *SANJBODController) GetDisks() []*SANDisk {
	return tjc.disksSlice
}

func (tjc *SANJBODController) GetWorkingDisk() *SANDisk {
	// todo more precise
	randomNum := randomInt(0, len(tjc.disks))
	return tjc.disks[randomNum]
}

func NewDiskState() *DiskState {
	return &DiskState{
		Async: 1,
	}
}
