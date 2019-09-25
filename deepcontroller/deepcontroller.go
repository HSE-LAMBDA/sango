package deepcontroller

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	lib "gosan"
	"log"
)

var _ = fmt.Println

type Manager struct {
	DeepControllingMode bool
	zmqReqSocket        *ZMQReqSocket
	sanComponents       map[string]DCAble
}

type DCAble interface {
	Break(*lib.Process, float64, float64)
	Repair(*lib.Process, float64)
	Update(map[string]float64)
	Reset()

	GetCurrentState() string
	SetCurrentState(string)
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

func (DCM *Manager) DeepControllerManagerProcess(TP *lib.Process, data interface{}) {
	args, ok := data.(*lib.DeepControllerFlags)
	if !ok {
		log.Panic("Deep controller conversion error")
	}
	if args.DeepControllingMode == false {
		log.Println("No need to communicate with DC")
		return
	}

	for {
		jsonBytes, err := json.MarshalIndent(DCM.sanComponents, "", "\t")
		if err != nil {
			log.Panic(err)
		}
		newParams := DCM.CommunicateWithDC(jsonBytes)
		DCM.CreateDeepControllerImpact(TP, newParams)

		for _, component := range DCM.sanComponents{
			component.Reset()
		}

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

func (DCM *Manager) Close() {
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

func NewManager(flags *lib.DeepControllerFlags) *Manager {
	dcm := &Manager{
		DeepControllingMode: flags.DeepControllingMode,
		sanComponents: make(map[string]DCAble),
	}

	dcm.zmqReqSocket = NewZMQReqSocket(flags.DeepProtocol, flags.DeepHost, flags.DeepPort)
	return dcm
}


func (DCM *Manager) CommunicateWithDC(jsonBytes []byte) *DeepParameters {
	//Send to deep controller
	parametersBytes := DCM.zmqReqSocket.request(jsonBytes)
	newParameters := new(DeepParameters)
	err := json.Unmarshal(parametersBytes, newParameters)
	if err != nil {
		log.Println(err)
	}
	return newParameters
}

func (DCM *Manager) CreateDeepControllerImpact(TP *lib.Process, newParams *DeepParameters) {
	DCM.ParametersUpdaterRecursiveDC(newParams.Components)
	DCM.AnomalyCreatorRecursiveDC(TP, newParams.Components)
}

func (DCM *Manager) ParametersUpdaterRecursiveDC(components map[string]*Component) {

	for name, component := range components {
		obj := DCM.getComponent(name)
		obj.Update(component.Parameters)
		DCM.ParametersUpdaterRecursiveDC(component.SubComponents)
	}
}

func (DCM *Manager) AnomalyCreatorRecursiveDC(TP *lib.Process, components map[string]*Component) {
	cTime := lib.SIM_get_clock()
	for name, component := range components {
		obj := DCM.getComponent(name)
		DCM.CreateAnomaly(TP, component.State, obj, cTime)
		DCM.AnomalyCreatorRecursiveDC(TP, component.SubComponents)
	}
}


func (DCM *Manager) getComponent(name string)DCAble{
	obj, ok := DCM.sanComponents[name]
	if !ok {
		log.Printf("Uknown component name %s has come from DC", name)
	}
	return obj
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