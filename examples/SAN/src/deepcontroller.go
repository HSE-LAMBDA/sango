package src

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"lib"
	"log"
)

var _ = fmt.Println

type DeepControllerManager struct {
	DeepControllingMode bool
	zmqReqSocket        *ZMQReqSocket
	updaterFuncMapDC    map[string]func(string, map[string]float64)
	anomalyFuncMapDC    map[string]func(*SANProcess, string, *ComponentAnomaly)
	clients             map[int]*ClientManager
	iob                 *IOBalancer
	tracer              *TracerManager
}

type (
	DeepParameters struct {
		Name       string                `json:"name"`
		Type       string                `json:"type"`
		Action     *int                  `json:"action"`
		Load       []*DeepLoad           `json:"load"`
		Components map[string]*Component `json:"subcomponents"`
	}
	DeepLoad struct {
		Name      string  `json:"name"`
		BlockSize string  `json:"blocksize"`
		Size      float64 `json:"size"`
		Time      float64 `json:"time"`
		Type      string  `json:"type"`
		NumJobs   int     `json:"num_jobs"`
	}

	Component struct {
		Type          string                `json:"type"`
		Name          string                `json:"name"`
		State         *ComponentAnomaly     `json:"state"`
		Parameters    map[string]float64    `json:"parameters"`
		SubComponents map[string]*Component `json:"subcomponents"`
	}

	ComponentAnomaly struct {
		Anomaly string  `json:"anomaly"`
		Degree  float64 `json:"degree"`
	}
)
type (
	ZMQReqSocket struct {
		host, protocol, port string
		requester            *zmq.Socket
	}
)

func (DCM *DeepControllerManager) DeepControllerManagerProcess(TP *SANProcess, data interface{}) {
	args, ok := data.(*DeepControllerFlags)
	if !ok {
		log.Panic("Deep controller conversion error")
	}
	if args.DeepControllingMode == false {
		log.Println("No need to communicate with DC")
		return
	}
	TP.Daemonize(DCM)

	tracer := DCM.tracer
	for {
		jsonBytes := tracer.PrepareJSONBytes()
		newParams := DCM.CommunicateWithDC(jsonBytes)
		DCM.CreateDeepControllerImpact(TP, newParams)
		tracer.ResetValues()
		TP.SIM_wait(args.DeepJsonDelay)
	}
}

func (h *ZMQReqSocket) request(logjson []byte) []byte {
	stateByte, _ := json.Marshal(logjson)
	_, err := h.requester.SendBytes(stateByte, 0)
	if err != nil {
		log.Panic("ZMQReqSocket error")
	}
	reply, _ := h.requester.RecvBytes(0)
	return reply
}

func NewZMQReqSocket(protocol string, host string, port string) *ZMQReqSocket {
	requester, _ := zmq.NewSocket(zmq.REQ)
	err := requester.Connect(protocol + "://" + host + ":" + port)
	if err != nil {
		log.Panic("ZMQReqSocket error")
	}
	return &ZMQReqSocket{
		host:      host,
		port:      port,
		protocol:  protocol,
		requester: requester,
	}
}

func (DCM *DeepControllerManager) Close() {
	socket := DCM.zmqReqSocket
	log.Println("Start sending DEAD to DC")

	//dead, _ := json.Marshal([]byte("DEAD"))
	//_, err := socket.SendBytes(dead, 0)
	//if err != nil {
	//	log.Panic("ZMQReqSocket error")
	//}
	a := socket.request([]byte("DEAD"))
	_ = a

	log.Println("Close ZMQ")
	err := DCM.zmqReqSocket.requester.Close()
	if err != nil {
		log.Panic(err)
	}
	return
}

func NewDeepControllerManager(flags *DeepControllerFlags, iob *IOBalancer, clients map[int]*ClientManager,
	t *TracerManager) *DeepControllerManager {
	dcm := &DeepControllerManager{
		DeepControllingMode: flags.DeepControllingMode,
		iob:                 iob,
		clients:             clients,
	}

	emptyFunc := func(*SANProcess, string, *ComponentAnomaly) {}
	anomalyFuncMapDC := make(map[string]func(*SANProcess, string, *ComponentAnomaly))
	updaterFuncMapDC := make(map[string]func(string, map[string]float64))

	anomalyFuncMapDC["controller"] = emptyFunc
	anomalyFuncMapDC["link"] = emptyFunc
	anomalyFuncMapDC["disk_drive"] = emptyFunc
	anomalyFuncMapDC["volume"] = emptyFunc
	anomalyFuncMapDC["cache"] = emptyFunc
	anomalyFuncMapDC["storage"] = emptyFunc
	anomalyFuncMapDC["pci-fabric"] = emptyFunc
	anomalyFuncMapDC["network_switch"] = emptyFunc
	anomalyFuncMapDC["io-balancer"] = emptyFunc

	updaterFuncMapDC["controller"] = dcm.ControllerUpdaterDC
	updaterFuncMapDC["link"] = dcm.LinkUpdaterDC
	updaterFuncMapDC["disk_drive"] = dcm.DiskUpdaterDC
	updaterFuncMapDC["storage"] = JBODUpdaterDC
	updaterFuncMapDC["pci-fabric"] = PCIeFabricUpdaterDC
	updaterFuncMapDC["network_switch"] = NetworkSwitchUpdaterDC
	updaterFuncMapDC["io-balancer"] = IOBalancerUpdaterDC
	updaterFuncMapDC["volume"] = VolumeUpdaterDC

	dcm.zmqReqSocket = NewZMQReqSocket(flags.DeepProtocol, flags.DeepHost, flags.DeepPort)
	dcm.anomalyFuncMapDC = anomalyFuncMapDC
	dcm.updaterFuncMapDC = updaterFuncMapDC
	dcm.tracer = t
	return dcm
}

func (DCM *DeepControllerManager) CreateLoad(newParams *DeepParameters) {
	load := newParams.Load
	clientHost := lib.GetHostByName("Client")

	for _, fileDC := range load {
		n := fileDC.NumJobs

		for _, client := range DCM.clients {
			if n > 0 {
				break
			}
			file := lib.NewFile(fileDC.Name, fileDC.Size, rwConversion(fileDC.Type), 1, fileDC.BlockSize, fileDC.Time)
			FORK("Client_"+file.Filename, client.PacketSenderReceiverProcess, clientHost, file)
			n--
		}
	}

	return
}

func (DCM *DeepControllerManager) CommunicateWithDC(jsonBytes []byte) *DeepParameters {
	//Send to deep controller
	parametersBytes := DCM.zmqReqSocket.request(jsonBytes)
	newParameters := new(DeepParameters)
	err := json.Unmarshal(parametersBytes, newParameters)
	if err != nil {
		log.Println(err)
	}
	return newParameters
}

func (DCM *DeepControllerManager) CreateDeepControllerImpact(TP *SANProcess, newParams *DeepParameters) {
	DCM.CreateLoad(newParams)
	DCM.ParametersUpdaterRecursiveDC(newParams.Components)
	DCM.AnomalyCreatorRecursiveDC(TP, newParams.Components)
}

func (DCM *DeepControllerManager) ParametersUpdaterRecursiveDC(components map[string]*Component) {

	for name, component := range components {
		function, ok := DCM.updaterFuncMapDC[component.Type]
		if !ok {
			log.Printf("Uknown component type %s", component.Type)
		}
		function(name, component.Parameters)
		DCM.ParametersUpdaterRecursiveDC(component.SubComponents)
	}
}

func (DCM *DeepControllerManager) AnomalyCreatorRecursiveDC(TP *SANProcess, components map[string]*Component) {
	cTime := lib.SIM_get_clock()
	for name, component := range components {
		breakComponent, ok := DCM.iob.allComponents[name]
		if !ok {
			log.Panicf("Uknown anomaly component %s", name)
		}

		DCM.CreateAnomaly(TP, component.State, breakComponent, cTime)
		DCM.AnomalyCreatorRecursiveDC(TP, component.SubComponents)
	}
}

func (DCM *DeepControllerManager) ControllerUpdaterDC(controllerId string, params map[string]float64) {
	iob := DCM.iob
	controller := iob.GetControllerById(controllerId)

	speed, ok := params["speed"]
	if !ok {
		log.Panic("No controller speed parameter from deepcontroller")
	}
	controller.Speed = speed
}

func (DCM *DeepControllerManager) DiskUpdaterDC(diskID string, params map[string]float64) {
	iob := DCM.iob
	disk := iob.GetDiskById(diskID)
	readRate, ok := params["read"]
	if !ok {
		log.Panic("No read disk parameter from deepcontroller")
	}
	disk.ReadLink.Bandwidth = readRate

	writeRate, ok := params["write"]
	if !ok {
		log.Panic("No write disk parameter from deepcontroller")
	}
	disk.WriteLink.Bandwidth = writeRate

	size, ok := params["size"]
	if !ok {
		log.Panic("No size disk parameter from deepcontroller")
	}
	disk.Size = size
}

func (DCM *DeepControllerManager) LinkUpdaterDC(linkID string, params map[string]float64) {
	iob := DCM.iob
	link := iob.GetLinkById(linkID)

	bandwidth, ok := params["bandwidth"]
	if !ok {
		log.Panic("No bandwidth parameter from deepcontroller")
	}
	link.Bandwidth = bandwidth

	latency, ok := params["latency"]
	if !ok {
		log.Panic("No latency parameter from deepcontroller")
	}
	link.Latency = latency
}

func JBODUpdaterDC(s string, float64s map[string]float64) {
	//todo nothing to update now
	return
}

func NetworkSwitchUpdaterDC(s string, float64s map[string]float64) {
	//todo nothing to update now
	return
}

func PCIeFabricUpdaterDC(s string, float64s map[string]float64) {
	//todo nothing to update now
	return
}

func IOBalancerUpdaterDC(s string, float64s map[string]float64) {
	//todo nothing to update now
	return
}

func VolumeUpdaterDC(s string, float64s map[string]float64) {
	//todo nothing to update now
	return
}

func JBODAndParametersUpdater(disks map[string]*lib.Storage, componentParameters map[string]*Component) {
	JBODs := lib.GetAllJBODs()

	// now there are no JBOD parameters to be updated
	_ = JBODs

	for name, disk := range disks {
		component, ok := componentParameters[name]
		if !ok {
			log.Panic(fmt.Sprintf("No %s parameters from deepcontroller", name))
		}
		params := component.Parameters

		readRate, ok := params["read"]
		if !ok {
			log.Panic("No read disk parameter from deepcontroller")
		}
		disk.ReadLink.Bandwidth = readRate

		writeRate, ok := params["write"]
		if !ok {
			log.Panic("No write disk parameter from deepcontroller")
		}
		disk.WriteLink.Bandwidth = writeRate

		size, ok := params["size"]
		if !ok {
			log.Panic("No size disk parameter from deepcontroller")
		}
		disk.Size = size
	}
}

var rwConversion func(t string) lib.RequestType

func rwConversionFactory() func(t string) lib.RequestType {
	rwMap := make(map[string]lib.RequestType)

	rwMap["r"] = lib.READ
	rwMap["w"] = lib.WRITE
	rwMap["randr"] = lib.RANDREAD
	rwMap["randw"] = lib.RANDWRITE

	return func(t string) lib.RequestType {
		xType, ok := rwMap[t]
		if !ok {
			log.Panicf("No such type of load %s", t)
		}
		return xType
	}
}
