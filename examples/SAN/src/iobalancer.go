//
// Created by kenenbek on 16.12.17.
//
package src

import (
	"fmt"
	"lib"
	"log"
	"math/rand"
)

var _ = fmt.Print

/*
=============================================
IOBalancer
=============================================
*/

type (
	IOBalancer struct {
		*NamingProps
		*IOBalancerProps `json:"props"`
		*lib.Host        `json:"-"`
		*SANComponent `json:"-"`
		GlobalDiskState  *DiskState                        `json:"-"`
		funcManager      *IOBalancerFManager               `json:"-"`
		Controllers      []*Controller                     `json:"-"`
		Jbods            []*SANJBODController           `json:"-"`
		controllersMap   map[string]*Controller            `json:"-"`
		linksMap         map[string]*SANLink            `json:"-"`
		jbodsMap         map[string]*SANJBODController  `json:"-"`
		disksMap         map[string]*SANDisk            `json:"-"`
		FailManager      *AtmosphereFailManager            `json:"-"`

		volumes       map[string]*SANVolume `json:"-"`
		openedFiles   map[string]*FileInfo     `json:"-"`
		packetCounter float64                  `json:"-"`
		allComponents map[string]BreakAble     `json:"-"`
	}
	IOBalancerFManager struct {
		controF SANBFunction `json:"-"`
		cacheWF SANBFunction `json:"-"`
		cacheRF SANBFunction `json:"-"`
		jbodWF  SANBFunction `json:"-"`
		jbodRF  SANBFunction `json:"-"`
	}
	IOBalancerProps struct {
		*CommonProps
		ReadResponseRate        float64 `json:"read_response_rate"`
		WriteResponseRate       float64 `json:"write_response_rate"`
		ReadRequestsRate        float64 `json:"read_requests_rate"`
		WriteRequestsRate       float64 `json:"write_requests_rate"`
		ReadDataVolume          float64 `json:"read_data_volume"`
		WriteDataVolume         float64 `json:"write_data_volume"`
		ReadRequestProcessTime  float64 `json:"read_request_process_time"`
		WriteRequestProcessTime float64 `json:"write_request_process_time"`
		IoProcessingMethod      uint8   `json:"io_processing_method"`
		ReadCancelRate          float64 `json:"read_cancel_rate"`
		WriteCancelRate         float64 `json:"write_cancel_rate"`
		BlockSize               uint16  `json:"block_size"`
		IOProcessCnt            uint8   `json:"io_process_cnt"`
		ReadQueueLength         uint8   `json:"read_queue_length"`
		WriteQueueLength        uint8   `json:"write_queue_length"`
	}
)

func (iob *IOBalancer) IOBalancerProcessManager(TP *SANProcess, data interface{}) {
	return
}

func NewIOBalancer(host *lib.Host, cons []*Controller, jbods []*SANJBODController) *IOBalancer {
	naming := NewNamingProps(host.Name, host.Type, host.Id)
	iob := &IOBalancer{
		NamingProps: naming,
		Host:        host,
		SANComponent: &SANComponent{
			currentState: "default",
		},
		IOBalancerProps: &IOBalancerProps{
			CommonProps: &CommonProps{
				Status: OK,
			},
		},
		Controllers: cons,
		Jbods:       jbods,

		controllersMap: make(map[string]*Controller),
		linksMap:       make(map[string]*SANLink),
		jbodsMap:       make(map[string]*SANJBODController),
		disksMap:       make(map[string]*SANDisk),

		openedFiles:   make(map[string]*FileInfo),
		allComponents: make(map[string]BreakAble),
	}
	// Need to initialize controllersMap
	for _, con := range cons {
		iob.controllersMap[con.NamingProps.Name] = con
		iob.allComponents[con.NamingProps.Name] = con
	}

	// Need to initialize disksMap
	for _, jbodCon := range jbods {
		iob.jbodsMap[jbodCon.NamingProps.Name] = jbodCon
		iob.allComponents[jbodCon.NamingProps.Name] = jbodCon

		for _, disk := range jbodCon.disks {
			iob.disksMap[disk.NamingProps.Name] = disk
			iob.allComponents[disk.NamingProps.Name] = disk
		}
	}

	iob.allComponents[iob.NamingProps.Name] = iob
	iob.allComponents["Client"] = nil

	// Need to initialize linksMap
	linksMap := lib.GetAllLinksMap()
	for name, link := range linksMap {
		SANLink := NewSANLink(link, iob)
		iob.linksMap[name] = SANLink
		iob.allComponents[name] = SANLink
	}

	return iob
}

// todo code duplication

func (iob *IOBalancer) SANComponentsPreparation(file *lib.File) (*FileInfo, *Controller, *SANJBODController) {
	controller, jbod := iob.GetWorkingComponents()
	fileInfo := NewFileInfo(file, controller, jbod)

	return fileInfo, controller, jbod
}

func (iob *IOBalancer) CreateSANExecutors(file *lib.File) *FileInfo {
	var fi *FileInfo
	switch file.RequestType {
	case lib.WRITE, lib.RANDWRITE:
		fi = iob.CreateSANWriterExecutors(file)
	case lib.READ, lib.RANDREAD:
		fi = iob.CreateSANReaderExecutors(file)
	default:
		log.Panic("No such request type")
	}
	return fi
}

func (iob *IOBalancer) CreateSANWriterExecutors(file *lib.File) *FileInfo {
	iob.WriteRequestsRate++

	fi, controller, jbod := iob.SANComponentsPreparation(file) // file info
	_ = jbod

	FORK(fi.clientContr, controller.ServerExecutor, controller.Host, fi)
	FORK(fi.cacheJbod, jbod.JbodWriterProcess, jbod.Host, fi)
	return fi
}

func (iob *IOBalancer) CreateSANReaderExecutors(file *lib.File) *FileInfo {
	iob.ReadRequestsRate++

	fi, controller, jbod := iob.SANComponentsPreparation(file) // file info

	FORK(fi.clientContr, controller.ServerExecutor, controller.Host, fi)
	FORK(fi.cacheJbod, jbod.JbodReaderProcess, jbod.Host, fi)
	return fi
}

func (iob *IOBalancer) GetWorkingComponents() (*Controller, *SANJBODController) {
	controller := iob.GetWorkingController()
	jbod := iob.GetWorkingJBOD()
	return controller, jbod
}

func (iob *IOBalancer) ResetDiffValues() {
	iob.ReadResponseRate = 0
	iob.WriteResponseRate = 0
	iob.ReadRequestProcessTime = 0
	iob.WriteRequestProcessTime = 0

	iob.packetCounter = 1
}

func (iob *IOBalancer) GetWorkingController() *Controller {
	if len(iob.Controllers) == 0 {
		log.Panic("No working controllers")
	}
	randomNumber := randomInt(0, len(iob.Controllers))
	return iob.Controllers[randomNumber]
}

func (iob *IOBalancer) GetWorkingJBOD() *SANJBODController {
	if len(iob.Jbods) == 0 {
		log.Panic("No working jbod")
	}
	randomNumber := randomInt(0, len(iob.Jbods))
	return iob.Jbods[randomNumber]
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func (CL *ClientManager) IOBalancerWriteCounterProcess(TP *SANProcess, _ interface{}) {
	TP.SetHost(CL.Host)
	testPacket := &lib.Packet{} //todo
	iob := CL.iob

	controller := iob.GetWorkingController()

	jbod := iob.GetWorkingJBOD()
	disk := jbod.GetWorkingDisk()

	t1 := lib.SIM_get_clock()
	status := TP.SendToHost(controller, testPacket)
	if status != lib.OK {
		return
	}
	TP.SetHost(controller.Host)

	t2 := lib.SIM_get_clock()
	status = controller.ExecutePacket(TP, testPacket, nil)
	if status != lib.OK {
		return
	}
	t3 := lib.SIM_get_clock()

	TP.WriteSync(disk.Storage, testPacket)

	t4 := lib.SIM_get_clock()

	iob.WriteResponseRate = 2 * (t4 - t1) // 2 is due to round trip
	iob.WriteRequestProcessTime = t3 - t2

	iob.GlobalDiskState.WriteAwaitMs = 2 * (t4 - t1)
	iob.GlobalDiskState.SvctmMs = t3 - t2

	return
}

func (CL *ClientManager) IOBalancerReadCounterProcess(TP *SANProcess, _ interface{}) {
	TP.SetHost(CL.Host)
	testPacket := &lib.Packet{} //todo
	iob := CL.iob

	controller := iob.GetWorkingController()
	jbod := iob.GetWorkingJBOD()
	disk := jbod.GetWorkingDisk()

	t1 := lib.SIM_get_clock()
	status := TP.SendToHost(controller, testPacket)
	if status != lib.OK {
		return
	}
	TP.SetHost(controller.Host)

	t2 := lib.SIM_get_clock()
	status = controller.ExecutePacket(TP, testPacket, nil)
	if status != lib.OK {
		return
	}
	t3 := lib.SIM_get_clock()

	TP.WriteSync(disk.Storage, testPacket)

	t4 := lib.SIM_get_clock()

	iob.ReadResponseRate = 2 * (t4 - t1) // 2 is due to round trip
	iob.ReadRequestProcessTime = t3 - t2

	iob.GlobalDiskState.ReadAwaitMs = 2 * (t4 - t1)
	iob.GlobalDiskState.SvctmMs = t3 - t2

	return
}

func (TP *SANProcess) SendToHost(controller *Controller, packet *lib.Packet) lib.STATUS {
	status := TP.SendToHostWithoutReceive(controller.Host, packet)
	return status
}

func (iob *IOBalancer) GetLinkById(linkId string) *SANLink {
	link, ok := iob.linksMap[linkId]
	if !ok {
		log.Panicf("No such SAN Link %s", linkId)
	}
	return link
}

func (iob *IOBalancer) GetControllerById(conId string) *Controller {
	con, ok := iob.controllersMap[conId]
	if !ok {
		log.Panicf("No such SAN controller %s", conId)
	}
	return con
}


func (iob *IOBalancer) GetJbodConById(jbodId string) *SANJBODController {
	jbod, ok := iob.jbodsMap[jbodId]
	if !ok {
		log.Panicf("No such SAN jbod %s", jbodId)
	}
	return jbod
}

func (iob *IOBalancer) GetDiskById(diskId string) *SANDisk {
	disk, ok := iob.disksMap[diskId]
	if !ok {
		log.Panicf("No such SAN disk %s", diskId)
	}
	return disk
}


func (iob *IOBalancer) RemoveWorkingController(index int) {
	iob.Controllers[index] = iob.Controllers[len(iob.Controllers)-1]
	iob.Controllers[len(iob.Controllers)-1] = nil
	iob.Controllers = iob.Controllers[:len(iob.Controllers)-1]
}

func (iob *IOBalancer) AddWorkingController(controller *Controller) {
	// add controller to the list of controllers
	index := len(iob.Controllers)
	controller.index = index
	iob.Controllers = append(iob.Controllers, controller)
}
