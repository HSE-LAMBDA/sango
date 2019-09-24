//
// Created by kenenbek on 18.12.17.
//

package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	lib "gosan"
	"log"
	"os"
)

var _ = fmt.Print

type (
	TracerManager struct {
		OutputFileWrite bool
		file            *os.File
		buffer          *bufio.Writer
		logs            *LogJson
		lastPollingTime float64
		iob             *IOBalancer
	}
)

func (tracer *TracerManager) TracerManagerProcess(TP *SANProcess, data interface{}) {
	args, ok := data.(*TracerFlags)
	if !ok {
		log.Panic("TracerFlags conversion error")
	}
	if args.outputFileName == "" {
		log.Println("No need to write output")
		return
	}

	TP.Daemonize(tracer)

	tracer.CreateOutputFile(args)

	for {

		tracer.PrepareJSONBytes()
		jsonBytes := tracer.PrepareJSONBytes()
		tracer.WriteToOutputFile(jsonBytes)
		//fmt.Println(tracer.logs.StState.KBWrtn)
		tracer.ResetValues()
		TP.SIM_wait(0.1)
	}
}

func NewTracerManager(iob *IOBalancer, acm *AtmosphereControlManager, cm *ClientManager) *TracerManager {
	arr := make([]TraceAble, len(iob.Controllers))
	for i, v := range iob.Controllers {
		arr[i] = v
	}

	hosts := append(arr, iob)
	tracer := &TracerManager{
		iob:  iob,
		logs: CreateLogJSON(hosts, iob.Jbods, acm),
	}
	return tracer
}

func (tracer *TracerManager) CreateOutputFile(flagsData *TracerFlags) {
	var err error
	log.Println("Output file will be generated")
	tracer.OutputFileWrite = true

	tracer.file, err = os.Create(flagsData.outputFileName)
	if err != nil {
		log.Panic(err)
	}
	tracer.buffer = bufio.NewWriter(tracer.file)
}

func (tracer *TracerManager) CloseOutputFile() {
	err := tracer.buffer.Flush()
	if err != nil {
		log.Panic(err)
	}
	err = tracer.file.Close()
	if err != nil {
		log.Panic(err)
	}
}

func (tracer *TracerManager) Close() {
	tracer.CloseOutputFile()
}

func CreateLogJSON(hosts []TraceAble, jbods []*SANJBODController, ACM *AtmosphereControlManager) *LogJson {
	/*
		hosts contains all iob, ns, pciefabric and controllers
	*/
	logjson := new(LogJson)

	logjson.Timestamp = lib.SIM_get_clock()
	logjson.Ambience = ACM.currentAtmosphere
	logjson.StState = NewDiskState()

	// Initialize logs at iobalancer, network switch, PCIe
	logjson.StorageComponents = append(logjson.StorageComponents, hosts...)

	for _, JBOD := range jbods {
		for _, disk := range JBOD.disks {
			logjson.StorageComponents = append(logjson.StorageComponents, disk)
		}
	}

	// add volumes
	iob, ok := hosts[len(hosts)-1].(*IOBalancer)
	if !ok {
		log.Panic("io-balancer conversion error")
	}
	for _, vol := range iob.volumes {
		logjson.StorageComponents = append(logjson.StorageComponents, vol)
	}
	return logjson
}

func (tracer *TracerManager) PrepareJSONBytes() []byte {
	tracer.UpdateLogJSON()
	jsonBytes, err := json.MarshalIndent(tracer.logs, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}

func (tracer *TracerManager) WriteToOutputFile(jsonBytes []byte) {
	_, err := tracer.buffer.Write(jsonBytes)
	if err != nil {
		log.Panic(err)
	}
}

func (tracer *TracerManager) UpdateLogJSON() {
	tracer.logs.Timestamp = lib.SIM_get_clock()
	tracer.AtmosphereTracer()
	tracer.ControllersTracer()
	tracer.IOBalancerTracer()
	tracer.DiskDriveTracer()
	tracer.VolumeTracer()

}

// Atmosphere tracer
func (tracer *TracerManager) AtmosphereTracer() {

}

func (tracer *TracerManager) ControllersTracer() {
	cons := tracer.iob.Controllers
	for _, controller := range cons {
		if controller.Status != FAIL {
			controller.Uptime += lib.SIM_get_clock() - tracer.lastPollingTime
		}
	}
}

//Disk tracer
func (tracer *TracerManager) DiskDriveTracer() {
	disks := tracer.iob.disksMap
	lastTime := tracer.lastPollingTime
	for _, disk := range disks {
		if disk.Status != FAIL {
			disk.Uptime += lib.SIM_get_clock() - lastTime
		}
	}
}

func (tracer *TracerManager) VolumeTracer() {
	for _, volume := range tracer.iob.volumes {
		lastTime := tracer.lastPollingTime

		capacity := 0.
		maxTemp := 0.0

		for _, disk := range volume.disks {
			if disk.Status != FAIL {
				disk.Uptime += lib.SIM_get_clock() - lastTime
				capacity += disk.UsedSpace
				if disk.DevTemp > maxTemp {
					maxTemp = disk.DevTemp
				}
			}
		}

		volume.UsedSpace = capacity
		volume.FreeSpace = volume.RawCapacity - capacity
		volume.DevTemp = maxTemp
	}
}

func (tracer *TracerManager) IOBalancerTracer() {
	iob := tracer.iob

	if iob.Status != FAIL {
		iob.Uptime += lib.SIM_get_clock() - tracer.lastPollingTime
	}
}

func (tracer *TracerManager) ResetValues() {
	components := tracer.logs.StorageComponents
	for _, component := range components {
		component.ResetDiffValues()
	}
	tracer.logs.StState.ResetDiffValues()
	tracer.lastPollingTime = lib.SIM_get_clock()
	return
}
