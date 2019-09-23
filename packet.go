package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type (
	RequestType int
)

var (
	//PACKET_4K    *Packet
	//PACKET_8K    *Packet
	//PACKET_16K   *Packet
	//PACKET_32K   *Packet
	//PACKET_64K   *Packet
	//PACKET_128K  *Packet
	//PACKET_256K  *Packet
	//PACKET_512K  *Packet
	//PACKET_1024K *Packet
	//PACKET_2048K *Packet
	PACKET_ACK      *Packet
	PACKET_FINALIZE *Packet
	PACKET_RECOVERY *Packet
)

const (
	_          = iota // ignore first value by assigning to blank identifier
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

const (
	RANDREAD RequestType = iota
	RANDWRITE
	READ
	WRITE
)

type File struct {
	Size        float64
	Filename    string
	RequestType RequestType
	Packet1     *Packet
	Packet2     *Packet
	Latency     float64
	NumJobs     int
}

type (
	CommonPacketInfo struct {
		Size  float64 `json:"size"`
		Index int     `json:"-"`
		// Network
		TransmissionTime    *NormalParam `json:"transmission_time"`
		LatencyTime         *NormalParam `json:"latency_time"`
		ReadProcessingTime  *NormalParam `json:"read_processing_time"`
		WriteProcessingTime *NormalParam `json:"write_processing_time"`
		RecoveryTime        *NormalParam `json:"recovery_time"`

		RandReadParams  *BlockParams `json:"rand_read"`
		RandWriteParams *BlockParams `json:"rand_write"`
		SeqReadParams   *BlockParams `json:"seq_read"`
		SeqWriteParams  *BlockParams `json:"seq_write"`
		ReadWriteParams *BlockParams `json:"read_write"`

		data            interface{}                  `json:"-"`
		allBlockParams  map[RequestType]*BlockParams `json:"-"`
		executionParams map[RequestType]*NormalParam `json:"-"`
	}

	Packet struct {
		Type  RequestType `json:"-"`
		Index int         `json:"-"`
		Size  float64     `json:"size"`

		NetworkTime *NetworkParams    `json:"transmission_time"`
		ContrTime   *ControllerParams `json:"recovery_time"`
		DiskParams  *BlockParams      `json:"disk_params"`
	}

	NetworkParams struct {
		TransmissionTime *NormalParam `json:"transmission_time"`
		LatencyTime      *NormalParam `json:"latency_time"`
	}

	ControllerParams struct {
		ProcessingTime *NormalParam `json:"processing_time"`
		RecoveryTime   *NormalParam `json:"recovery_time"`
	}

	BlockParams struct {
		// Storage
		RateTime      *NormalParam `json:"rate_time"`
		SeekTime      *NormalParam `json:"seek_time"`
		OverheadsTime *NormalParam `json:"overheads_time"`
	}

	NormalParam struct {
		Mean   float64 `json:"mean"`
		StdDev float64 `json:"std_dev"`
	}
)

func InitPacketInfo(filename string) {
	if filename == "" {
		log.Println("Working without packet mode")
		return
	}

	var PACKET_INFO map[string]*CommonPacketInfo
	packet_json, err := os.Open(filename)

	packets := make([]*CommonPacketInfo, 10)
	packetIndices := make([]string, 10)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer packet_json.Close()

	byteValue, _ := ioutil.ReadAll(packet_json)

	jerr := json.Unmarshal(byteValue, &PACKET_INFO)
	if jerr != nil {
		log.Panic(jerr)
	}

	for i := 0; i < 10; i++ {
		num := int(math.Pow(2, float64(2+i)))
		name := strconv.Itoa(num)
		packetIndices[i] = name + "K"
	}

	for i, packetName := range packetIndices {
		genPacket, ok := PACKET_INFO[packetName]
		if !ok {
			log.Panic("no such packet name")
		}
		packets[i] = genPacket
	}

	for i := 0; i < len(packets); i++ {
		packets[i].Index = i

		packets[i].allBlockParams = make(map[RequestType]*BlockParams)
		packets[i].executionParams = make(map[RequestType]*NormalParam)

		packets[i].allBlockParams[RANDREAD] = packets[i].RandReadParams
		packets[i].allBlockParams[RANDWRITE] = packets[i].RandWriteParams
		packets[i].allBlockParams[READ] = packets[i].SeqReadParams
		packets[i].allBlockParams[WRITE] = packets[i].SeqWriteParams

		packets[i].executionParams[RANDREAD] = packets[i].ReadProcessingTime
		packets[i].executionParams[RANDWRITE] = packets[i].WriteProcessingTime
		packets[i].executionParams[READ] = packets[i].ReadProcessingTime
		packets[i].executionParams[WRITE] = packets[i].WriteProcessingTime
	}

	concretePackets := make(map[RequestType]map[string]*Packet)
	concretePackets[RANDREAD] = make(map[string]*Packet)
	concretePackets[RANDWRITE] = make(map[string]*Packet)
	concretePackets[READ] = make(map[string]*Packet)
	concretePackets[WRITE] = make(map[string]*Packet)

	for reqType, packetsMap := range concretePackets {
		for packetName, commonPacketInfo := range PACKET_INFO {
			packet := newPacket(commonPacketInfo, reqType)
			packetsMap[packetName] = packet
		}
	}

	normalParamMockup := &NormalParam{
		Mean:   1e-6,
		StdDev: 0,
	}

	PACKET_ACK = &Packet{
		Size: 1024,
		NetworkTime: &NetworkParams{
			TransmissionTime: normalParamMockup,
			LatencyTime: normalParamMockup,
		},
		ContrTime: &ControllerParams{
			ProcessingTime: normalParamMockup,
			RecoveryTime: normalParamMockup,
		},
	}
	PACKET_FINALIZE = &Packet{
		Size: 0,
		NetworkTime: &NetworkParams{
			TransmissionTime: &NormalParam{
				Mean:   1e-6,
				StdDev: 0,
			},
			LatencyTime: &NormalParam{
				Mean:   1e-6,
				StdDev: 0,
			},
		},
	}

	initPacketClosures(concretePackets, packetIndices)
}

func NewFile(filename string, size float64, requestType RequestType, numJobs int, blockSize string,
	latency float64) *File {

	var p1, p2 *Packet
	packet := GetPacketByName(requestType, blockSize)

	switch requestType {
	case RANDWRITE, WRITE:
		p1 = packet
		p2 = PACKET_ACK
	case RANDREAD, READ:
		p1 = PACKET_ACK
		p2 = packet
	default:
		// lol kek cheburek
	}

	f := &File{
		Size:        size,
		Filename:    filename,
		RequestType: requestType,
		Packet1:     p1,
		Packet2:     p2,
		Latency:     latency,
		NumJobs:     numJobs,
	}
	return f
}

func newPacket(cpi *CommonPacketInfo, rType RequestType) *Packet {
	diskParams := getBlockParametersInfo(cpi, rType)
	execParams := getExecutionParametersInfo(cpi, rType)

	return &Packet{
		Type:  rType,
		Size:  cpi.Size,
		Index: cpi.Index,

		NetworkTime: &NetworkParams{
			TransmissionTime: cpi.TransmissionTime,
			LatencyTime:      cpi.LatencyTime,
		},
		ContrTime: &ControllerParams{
			ProcessingTime: execParams,
			RecoveryTime:   cpi.RecoveryTime,
		},
		DiskParams: diskParams,
	}
}

func getRandomPacketFactory(cpi map[RequestType]map[string]*Packet, packets []string) (func(RequestType) *Packet,
	func() string) {
	f1 := func(requestType RequestType) *Packet {
		randomValue := rand.Intn(10)
		packetMap, ok := cpi[requestType]
		if !ok {
			log.Panic("No such request type")
		}
		randomPacketName := packets[randomValue]
		packet, ok := packetMap[randomPacketName]
		if !ok {
			log.Panic("No such packet")
		}
		return packet
	}
	f2 := func() string {
		randomValue := rand.Intn(10)
		return packets[randomValue]
	}
	return f1, f2
}

func getPacketByNameFactory(cpInfo map[RequestType]map[string]*Packet) func(RequestType, string) *Packet {
	return func(requestType RequestType, packetName string) *Packet {
		packetMap, ok := cpInfo[requestType]
		if !ok {
			log.Panic("No such request type")
		}
		packet, ok := packetMap[packetName]
		if !ok {
			log.Panic("No such packet")
		}
		return packet
	}
}
