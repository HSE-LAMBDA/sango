package lib

import (
	"log"
	"math/rand"
)

type Resource interface {
	Put(e *Event, globalQueue *globalEventQueue)
	//getState()float64
	//getLastTime()float64
}

func getTransferTime(resource Resource, packet *Packet) (float64, float64) {
	link, ok := resource.(*Link)
	if !ok {
		log.Panic("Wrong link conversion")
	}

	t1 := Maximum(SIM_get_clock(), link.lastTime)
	tt := packet.NetworkTime.TransmissionTime
	lt := packet.NetworkTime.LatencyTime

	sendTime := link.Bandwidth * (rand.NormFloat64()*tt.StdDev + tt.Mean)
	latencyTime := rand.NormFloat64()*lt.StdDev + lt.Mean

	timeEnd := t1 + sendTime + latencyTime

	return timeEnd, sendTime + latencyTime
}

func getExecutionTime(resource Resource, packet *Packet) (float64, float64) {
	core, ok := resource.(*Core)
	if !ok {
		log.Panic("Wrong core conversion")
	}
	et := packet.ContrTime.ProcessingTime

	t1 := Maximum(SIM_get_clock(), core.lastTime)
	execTime := (rand.NormFloat64()*et.StdDev + et.Mean) / (core.speed * core.state)

	timeEnd := t1 + execTime
	return timeEnd, execTime
}

func getRecoveryTime(resource Resource, packet *Packet) (float64, float64) {
	core, ok := resource.(*Core)
	if !ok {
		log.Panic("Wrong core conversion")
	}
	rt := packet.ContrTime.RecoveryTime

	t1 := Maximum(SIM_get_clock(), core.lastTime)
	rTime := (rand.NormFloat64()*rt.StdDev + rt.Mean) / (core.speed * core.state)

	timeEnd := t1 + rTime
	return timeEnd, rTime
}

func getDiskOperationTime(resource Resource, packet *Packet) (float64, float64) {
	link, ok := resource.(*Link)
	if !ok {
		log.Panic("Wrong link conversion")
	}
	params := packet.DiskParams

	t1 := Maximum(SIM_get_clock(), link.lastTime)
	writeTime := link.Bandwidth * (rand.NormFloat64()*params.RateTime.StdDev + params.RateTime.Mean)
	seekTime := rand.NormFloat64()*params.SeekTime.StdDev + params.SeekTime.Mean
	overheadTime := rand.NormFloat64()*params.OverheadsTime.StdDev + params.OverheadsTime.Mean

	timeEnd := t1 + writeTime + seekTime + overheadTime

	return timeEnd, writeTime + seekTime + overheadTime
}

func getBlockParametersInfo(packet *CommonPacketInfo, pType RequestType) *BlockParams {
	params, ok := packet.allBlockParams[pType]
	if !ok {
		log.Panic("not such packet Type")
	}
	return params
}

func getExecutionParametersInfo(packet *CommonPacketInfo, pType RequestType) *NormalParam {
	params, ok := packet.executionParams[pType]
	if !ok {
		log.Panic("not such packet Type")
	}
	return params
}
