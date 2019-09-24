//
// Created by kenenbek on 18.12.17.
//
package src

import (
	"fmt"
	lib "gosan"
)

var _ = fmt.Print

type (
	SANDisk struct {
		*NamingProps
		*StorageProperties `json:"props"`
		*lib.Storage       `json:"-"`
		*SANComponent      `json:"-"`
		GlobalDiskState    *DiskState         `json:"-"`
		iob                *IOBalancer        `json:"-"`
		jbod               *SANJBODController `json:"-"`
		index              int                `json:"-"`
	}

	StorageProperties struct {
		*CommonProps
		RawCapacity       float64 `json:"cap_gb"`
		Type              string  `json:"type"`
		JBOD              string  `json:"jbod"`
		Pool              string  `json:"pool"`
		AvgReadSpeed      float64 `json:"-"` // `json:"avg_read_speed"`
		AvgWriteSpeed     float64 `json:"-"` // `json:"avg_write_speed"`
		DataInterfacesCnt uint16  `json:"-"` // `json:"data_interfaces_cnt"`
		UsedSpace         float64 `json:"-"` // `json:"used_space"`
		FreeSpace         float64 `json:"-"` // `json:"free_space"`
	}
)

func (disk *SANDisk) DiskWriterManagerProcess(TP *SANProcess, data interface{}) {
	TP.Daemonize()
	hostName := TP.GetHost().GetName()
	mailbox := hostName + "_"

	for {
		task, _ := TP.ReceivePacket(mailbox)
		disk.WriteAsync(TP, task)
	}
}

func NewSANDisk(libDisk *lib.Storage, jbod *SANJBODController) *SANDisk {
	naming := NewNamingProps(libDisk.Name, libDisk.Type, libDisk.ID)
	td := &SANDisk{
		NamingProps: naming,
		StorageProperties: &StorageProperties{

			FreeSpace:     libDisk.Size,
			RawCapacity:   libDisk.Size,
			AvgReadSpeed:  libDisk.ReadRate,
			AvgWriteSpeed: libDisk.WriteRate,

			CommonProps: &CommonProps{
				Status: OK,
			},
		},
		SANComponent: &SANComponent{
			currentState: "default",
		},
		Storage: libDisk,
		jbod:    jbod,
	}
	return td
}

func (ds *DiskState) ResetDiffValues() {
	ds.ReadOps = 0
	ds.WriteOps = 0
	ds.ReadAwaitMs = 0
	ds.WriteAwaitMs = 0
	ds.Async = 0
	ds.Crrq = 0
	ds.AvgQsz = 0
	ds.SvctmMs = 0
	ds.MBRead = 0
	ds.MBWrtn = 0

	for i := 0; i < len(ds.IORead); i++ {
		ds.IORead[i] = 0
	}

	for i := 0; i < len(ds.IOWrite); i++ {
		ds.IOWrite[i] = 0
	}

}

func (disk *SANDisk) ResetDiffValues() {

}

func (volume *SANVolume) ResetDiffValues() {

}

func (disk *SANDisk) WriteAsync(TP *SANProcess, packet *lib.Packet) lib.STATUS {
	s := TP.WriteAsync(disk.Storage, packet)
	disk.MetricWrite(packet)
	return s

}

func (disk *SANDisk) WriteSync(TP *SANProcess, packet *lib.Packet) lib.STATUS {
	s := TP.WriteSync(disk.Storage, packet)
	disk.MetricWrite(packet)
	return s
}

func (disk *SANDisk) MetricWrite(packet *lib.Packet) {
	disk.UsedSpace += packet.Size
	disk.FreeSpace -= packet.Size

	disk.iob.WriteDataVolume += packet.Size

	disk.GlobalDiskState.MBWrtn += float64(packet.Size) / 1024
	disk.GlobalDiskState.WriteOps++
	disk.GlobalDiskState.IOWrite[packet.Index]++
}

func (disk *SANDisk) ReadAsync(TP *SANProcess, packet *lib.Packet) (*lib.Packet, lib.STATUS) {
	p, s := TP.ReadAsync(disk.Storage, packet)
	disk.MetricRead(packet)
	return p, s
}

func (disk *SANDisk) ReadSync(TP *SANProcess, packet *lib.Packet) (*lib.Packet, lib.STATUS) {
	p, s := TP.ReadSync(disk.Storage, packet)
	disk.MetricRead(packet)
	return p, s
}

func (disk *SANDisk) MetricRead(packet *lib.Packet) {
	disk.iob.ReadDataVolume += packet.Size

	disk.GlobalDiskState.ReadOps++
	disk.GlobalDiskState.MBRead += float64(packet.Size) / 1024
	disk.GlobalDiskState.IORead[packet.Index]++
}

func DiskReaderExecutor(TP *SANProcess, data interface{}) {
	//data := p.GetData().(*FileInfo)
	//
	//receiveAddress1 := data.diskReaderExecutor
	//sendAddress1 := data.cacheExecutorACK
	//
	//readDataTask := lib.NewTask("", 0, packetSize, nil)
	//for {
	//	task, res := p.ReceivePacket(receiveAddress1)
	//	if res != lib.OK || strings.Compare(task.GetName(), "finalize") == 0 {
	//		return
	//	}
	//	res = p.SendPacket(readDataTask, sendAddress1)
	//	if res != lib.OK {
	//		return
	//	}
	//}
}
