//
// Created by kenenbek on 29.11.17.
//

package src

import (
	"fmt"
	lib "gosan"
	"log"
)

var _ = fmt.Print

type (
	Controller struct {
		*NamingProps
		*ControllerProps `json:"props"`
		*lib.Host        `json:"-"`
		*SANComponent    `json:"-"`
		index            int         `json:"-"`
		iob              *IOBalancer `json:"-"`
	}
	ControllerProps struct {
		*CommonProps
		LastTime float64 `json:"-"`

		TrafficDif       float64 `json:"traffic"`
		TrafficCum       float64 `json:"-"`
		InputTrafficDif  float64 `json:"-"`
		OutputTrafficDif float64 `json:"-"`
		InputTrafficCum  float64 `json:"-"`
		OutputTrafficCum float64 `json:"-"`

		FlopsExecutedDiff float64 `json:"load"`
		FlopsExecutedCum  float64 `json:"-"`
	}
)

func (controller *Controller) ServerExecutor(TP *SANProcess, data interface{}) {
	arg, ok := data.(*FileInfo)
	if !ok {
		log.Panic("File info conversion error")
	}

	receiveAddress1 := arg.clientContr
	sendAddress1 := arg.cacheJbod
	//sendAddress1 := arg.contrCache

	receiveAddress2 := arg.jbodCache
	//receiveAddress2 := arg.cacheContr
	sendAddress2 := arg.contrClient

	for {
		packet, res := controller.ReceivePacket(TP, receiveAddress1)

		if res != lib.OK || packet == lib.PACKET_FINALIZE {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress1)
			return
		}

		res = controller.ExecutePacket(TP, packet, arg)
		if res != lib.OK {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress1)
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress2)
			return
		}

		// Send data to Fabric Manager
		res = controller.SendPacket(TP, packet, sendAddress1)
		if res != lib.OK {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress2)
			return
		}
		rTask, res := controller.ReceivePacket(TP, receiveAddress2)
		if res != lib.OK || packet == lib.PACKET_FINALIZE {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress2)
			return
		}

		res = controller.SendPacket(TP, rTask, sendAddress2)
		if res != lib.OK {
			TP.SendPacket(lib.PACKET_FINALIZE, sendAddress1)
			return
		}
	}
}

func NewController(host *lib.Host) *Controller {
	naming := NewNamingProps(host.Name, host.Type, host.Id)
	controller := &Controller{
		NamingProps: naming,
		Host:        host,
		SANComponent: &SANComponent{
			currentState: "default",
		},
		ControllerProps: &ControllerProps{
			CommonProps: &CommonProps{
				Status: OK,
			},
		},
	}
	return controller
}

func (controller *Controller) SendPacket(TP *SANProcess, packet *lib.Packet, address string) lib.STATUS {
	s := TP.SendPacket(packet, address)

	controller.OutputTrafficDif += packet.Size
	controller.OutputTrafficCum += packet.Size
	controller.TrafficDif += packet.Size
	controller.TrafficCum += packet.Size
	return s
}

func (controller *Controller) ReceivePacket(TP *SANProcess, address string) (*lib.Packet, lib.STATUS) {
	packet, s := TP.ReceivePacket(address)

	controller.InputTrafficDif += packet.Size
	controller.InputTrafficCum += packet.Size
	controller.TrafficDif += packet.Size
	controller.TrafficCum += packet.Size
	return packet, s
}

func (controller *Controller) ExecutePacket(TP *SANProcess, packet *lib.Packet, fi *FileInfo) lib.STATUS {
	iob := controller.iob

	t1 := lib.SIM_get_clock()
	s, flops := TP.Execute(packet)
	t2 := lib.SIM_get_clock()

	controller.FlopsExecutedDiff += flops
	controller.FlopsExecutedCum += flops

	deltaTime := t2 - t1
	switch fi.RequestType {
	case lib.RANDWRITE, lib.WRITE:
		t := iob.WriteRequestProcessTime + (deltaTime-iob.WriteRequestProcessTime)/iob.packetCounter
		iob.WriteRequestProcessTime = t
		iob.GlobalDiskState.SvctmMs = t
	case lib.RANDREAD, lib.READ:
		t := iob.ReadRequestProcessTime + (deltaTime-iob.ReadRequestProcessTime)/iob.packetCounter
		iob.ReadRequestProcessTime = t
		iob.GlobalDiskState.SvctmMs = t
	}

	return s
}

func (controller *Controller) ResetDiffValues() {
	controller.InputTrafficDif = 0
	controller.OutputTrafficDif = 0
	controller.TrafficDif = 0
	controller.FlopsExecutedDiff = 0
}
