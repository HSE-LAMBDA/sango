package lib

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
	Parse platform file
*/

type topology struct {
	XMLName       xml.Name           `xml:"topology"`
	ID            string             `xml:"id,attr"`
	Storage_types []storage_type_xml `xml:"storage_type"`
	Hosts         []host_xml         `xml:"host"`

	Iobalancer    iobalancer_xml     `xml:"iobalancer"`
	FabricManager fabric_manager_xml `xml:"fabric_manager"`
	NetworkSwitch network_switch_xml `xml:"network_switch"`

	Links        []link_xml         `xml:"link"`
	Routes       []route_xml        `xml:"route"`
	BackupRoutes []backup_route_xml `xml:"backup_route"`
	JBODS        []jbod_xml         `xml:"jbod"`
	CacheJbods   []cache_xml        `xml:"cache"`

	Volumes []*Volume `xml:"volume"`
}

type storage_type_xml struct {
	XMLName     xml.Name         `xml:"storage_type"`
	ID          string           `xml:"id,attr"`
	Size        string           `xml:"size,attr"`
	Model_props []model_prop_xml `xml:"model_prop"`
}

type model_prop_xml struct {
	XMLName xml.Name `xml:"model_prop"`
	Id      string   `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
}

type jbod_xml struct {
	XMLName       xml.Name `xml:"jbod"`
	Id            string   `xml:"id,attr"`
	StorageTypeId string   `xml:"storage_id,attr"`
	Amount        int      `xml:"amount,attr"`
}

type cache_xml struct {
	XMLName       xml.Name `xml:"cache"`
	Id            string   `xml:"id,attr"`
	StorageTypeId string   `xml:"storage_id,attr"`
	Amount        int      `xml:"amount,attr"`
}

type basic_host_xml struct {
	Id     string      `xml:"id,attr"`
	Speed  string      `xml:"speed,attr"`
	Mounts []mount_xml `xml:"mount"`
}

type host_xml struct {
	XMLName xml.Name    `xml:"host"`
	Id      string      `xml:"id,attr"`
	Type    string      `xml:"type,attr"`
	Speed   string      `xml:"speed,attr"`
	NCore   string      `xml:"nCore,attr"`
	Mounts  []mount_xml `xml:"mount"`
	//*basic_host_xml
}

type iobalancer_xml struct {
	XMLName xml.Name    `xml:"iobalancer"`
	Id      string      `xml:"id,attr"`
	Speed   string      `xml:"speed,attr"`
	Mounts  []mount_xml `xml:"mount"`
	//*basic_host_xml
}

type network_switch_xml struct {
	XMLName xml.Name    `xml:"network_switch"`
	Id      string      `xml:"id,attr"`
	Speed   string      `xml:"speed,attr"`
	Mounts  []mount_xml `xml:"mount"`
	//*basic_host_xml
}

type fabric_manager_xml struct {
	XMLName xml.Name    `xml:"fabric_manager"`
	Id      string      `xml:"id,attr"`
	Speed   string      `xml:"speed,attr"`
	Mounts  []mount_xml `xml:"mount"`
	//*basic_host_xml
}

type mount_xml struct {
	XMLName   xml.Name `xml:"mount"`
	StorageId string   `xml:"storageId,attr"`
	Name      string   `xml:"name,attr"`
}

type link_xml struct {
	XMLName   xml.Name `xml:"link"`
	Id        string   `xml:"id,attr"`
	Bandwidth string   `xml:"bandwidth,attr"`
	Latency   string   `xml:"latency,attr"`
}

type route_xml struct {
	XMLName   xml.Name       `xml:"route"`
	Src       string         `xml:"src,attr"`
	Dst       string         `xml:"dst,attr"`
	Link_ctns []link_ctn_xml `xml:"link_ctn"`
}

type backup_route_xml struct {
	XMLName   xml.Name       `xml:"backup_route"`
	Src       string         `xml:"src,attr"`
	Dst       string         `xml:"dst,attr"`
	Link_ctns []link_ctn_xml `xml:"link_ctn"`
}

type link_ctn_xml struct {
	XMLName xml.Name `xml:"link_ctn"`
	Id      string   `xml:"id,attr"`
}

func SIM_platform_init(FilePath string, flag_disk_amount int) {
	// Open our xmlFile
	xmlFile, err := os.Open(FilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Panic(err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var TOPOLOGY topology
	xerr := xml.Unmarshal(byteValue, &TOPOLOGY)
	if xerr != nil {
		log.Panic(xerr)
	}

	JBODMap := make(map[string]*JBOD)
	JBODSlice := make([]*JBOD, 0)

	hostsMap := make(map[string]*Host)
	hostLinksMap := make(map[*Host][]*Link)
	hostHostLinkMap := make(map[*Host]map[*Host]*Link)
	hostHostLinkMapB := make(map[*Host]map[*Host]*Link)

	// STORAGE_TYPES
	StorageTypes := make(map[string]*StorageType)
	for i := 0; i < len(TOPOLOGY.Storage_types); i++ {
		s := StorageType{
			TypeId: TOPOLOGY.Storage_types[i].ID,
			Size:   UnitToFloatParser(TOPOLOGY.Storage_types[i].Size),
		}
		for j := 0; j < len(TOPOLOGY.Storage_types[i].Model_props); j++ {
			speed := TOPOLOGY.Storage_types[i].Model_props[j].Id
			value := TOPOLOGY.Storage_types[i].Model_props[j].Value
			if strings.Compare(speed, "read") == 0 {
				s.ReadRate = UnitToFloatParser(value)
			} else if strings.Compare(speed, "write") == 0 {
				s.WriteRate = UnitToFloatParser(value)
			}
		}
		StorageTypes[TOPOLOGY.Storage_types[i].ID] = &s
	}

	// Storage arrays
	var diskAmount int
	for i := 0; i < len(TOPOLOGY.JBODS); i++ {
		JBODId := TOPOLOGY.JBODS[i].Id
		if flag_disk_amount == -1 {
			diskAmount = TOPOLOGY.JBODS[i].Amount
		}
		jbod := NewJBOD(JBODId, diskAmount, TOPOLOGY.JBODS[i].StorageTypeId, StorageTypes)
		JBODMap[JBODId] = jbod
		JBODSlice = append(JBODSlice, jbod)

	}

	// Cache jbods
	for i := 0; i < len(TOPOLOGY.CacheJbods); i++ {
		cacheJbodId := TOPOLOGY.CacheJbods[i].Id
		cacheJbod := NewJBOD(cacheJbodId, TOPOLOGY.CacheJbods[i].Amount, TOPOLOGY.CacheJbods[i].StorageTypeId, StorageTypes)
		initGetCacheJbodClosure(cacheJbod) // todo
	}

	// HOSTS
	for i := 0; i < len(TOPOLOGY.Hosts); i++ {
		HostName := TOPOLOGY.Hosts[i].Id
		HostType := TOPOLOGY.Hosts[i].Type
		HostSpeed := UnitToFloatParser(TOPOLOGY.Hosts[i].Speed)
		HostCores, _ := strconv.Atoi(TOPOLOGY.Hosts[i].NCore)
		host := NewHost(HostName, HostType, HostName, HostSpeed, uint64(HostCores))
		hostsMap[HostName] = host

		hostHostLinkMap[host] = make(map[*Host]*Link)
	}

	initHostClosures(hostsMap)

	// Volumes
	initGetVolumeClosure(TOPOLOGY.Volumes)

	// LINKS
	MapLink := make(map[string]*Link)
	for i := 0; i < len(TOPOLOGY.Links); i++ {
		RealLinkId := TOPOLOGY.Links[i].Id
		RealLinkBW := UnitToFloatParser(TOPOLOGY.Links[i].Bandwidth)

		link := NewLink(RealLinkBW, RealLinkId)
		MapLink[RealLinkId] = link
	}

	// ROUTES
	for i := 0; i < len(TOPOLOGY.Routes); i++ {
		SRCHost := GetHostByName(TOPOLOGY.Routes[i].Src)
		DSTHost := GetHostByName(TOPOLOGY.Routes[i].Dst)
		RealLinkId := TOPOLOGY.Routes[i].Link_ctns[0].Id
		link, ok := MapLink[RealLinkId]
		if !ok {
			log.Panic("No such link in platform file")
		}
		link.Src = SRCHost
		link.Dst = DSTHost
		hostHostLinkMap[SRCHost][DSTHost] = link
		hostHostLinkMap[DSTHost][SRCHost] = link
	}

	// Backup routes
	/*for i := 0; i < len(TOPOLOGY.BackupRoutes); i++ {
		SRCHost := env.getHostByName(TOPOLOGY.BackupRoutes[i].Src).(*Host)
		DSTHost := env.getHostByName(TOPOLOGY.BackupRoutes[i].Dst).(*Host)

		links := make([]*Link, len(TOPOLOGY.BackupRoutes[i].Link_ctns))
		reverse_links := make([]*Link, len(TOPOLOGY.BackupRoutes[i].Link_ctns))
		for j := 0; j < len(TOPOLOGY.BackupRoutes[i].Link_ctns); j++ {
			RealLinkId := TOPOLOGY.BackupRoutes[i].Link_ctns[j].Id
			links[j] = MapLink[RealLinkId]
			reverse_links[len(reverse_links)-j-1] = MapLink[RealLinkId]
		}
		RealRoute := Route{start: SRCHost, finish: DSTHost}
		ReverseRealRoute := Route{start: DSTHost, finish: SRCHost}
		backupRoutes[RealRoute] = links
		backupRoutes[ReverseRealRoute] = reverse_links
	}*/

	// Find All income and outcome links for each host
	for srcHost, hostMap := range hostHostLinkMap {
		for _, link := range hostMap {
			hostLinksMap[srcHost] = append(hostLinksMap[srcHost], link)
		}
	}

	initLinkMapClosures(MapLink, hostHostLinkMap, hostHostLinkMapB)
	initHostLinkClosures(hostLinksMap)
	initJBODClosures(JBODMap, JBODSlice)
}

/*
	Parse deployment file
*/

type deployment struct {
	XMLName   xml.Name    `xml:"deployment"`
	Processes []process_x `xml:"process"`
}

type process_x struct {
	XMLName   xml.Name   `xml:"process"`
	Host      string     `xml:"host,attr"`
	Function  string     `xml:"function,attr"`
	Arguments []argument `xml:"argument"`
}

type argument struct {
	XMLName xml.Name `xml:"argument"`
	Value   string   `xml:"value,attr"`
}

func SIM_launch_application(FilePath string) {
	// Open our xmlFile
	xmlFile, err := os.Open(FilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var DEPLOYMENT deployment
	xerr := xml.Unmarshal(byteValue, &DEPLOYMENT)
	if xerr != nil {
		log.Panic(xerr)
	}

	fMap := make(map[string]BaseBehaviourFunction)
	fData := make(map[string]interface{})
	initDeploymentClosure(fMap, fData)

	for i := 0; i < len(DEPLOYMENT.Processes); i++ {
		HostName := DEPLOYMENT.Processes[i].Host
		FuncName := DEPLOYMENT.Processes[i].Function
		Arguments := argsToStr(DEPLOYMENT.Processes[i].Arguments)
		Func := GetFunctionByName(FuncName)
		Host := GetHostByName(HostName)

		SIM_subprocess_create_async(FuncName, Func, Host, nil, Arguments)
	}
}

func (TOPOLOGY *topology) getLinkById(ID string) link_xml {
	for i := 0; i < len(TOPOLOGY.Links); i++ {
		if TOPOLOGY.Links[i].Id == ID {
			return TOPOLOGY.Links[i]
		}
	}
	panic("No such link id")
	return link_xml{}
}

func getFunctionByNameFactory(FunctionsMap map[string]BaseBehaviourFunction) func(string) BaseBehaviourFunction {
	return func(FuncName string) BaseBehaviourFunction {
		Func, ok := FunctionsMap[FuncName]
		if !ok {
			panic(fmt.Sprintf("No such registered function %s", FuncName))
		}
		return Func
	}
}

func argsToStr(args []argument) []string {
	array := make([]string, len(args))
	for i := range args {
		array[i] = args[i].Value
	}
	return array
}

func UnitToFloatFactory() func(string) float64 {
	unitsMap := make(map[string]float64)
	TERA := math.Pow10(12)
	GIGA := math.Pow10(9)
	MEGA := math.Pow10(6)
	KILO := math.Pow10(3)

	unitsMap["TB"] = TERA
	unitsMap["GB"] = GIGA
	unitsMap["MB"] = MEGA
	unitsMap["KB"] = KILO
	unitsMap["B"] = 1

	unitsMap["GBps"] = GIGA
	unitsMap["MBps"] = MEGA
	unitsMap["KBps"] = KILO
	unitsMap["Bps"] = 1

	unitsMap["Gf"] = GIGA
	unitsMap["Mf"] = MEGA
	unitsMap["Kf"] = KILO
	unitsMap["f"] = 1

	unitsMap["s"] = 1
	unitsMap["ms"] = 0.001
	unitsMap["us"] = 0.000001

	re := regexp.MustCompile("\\d*\\.?\\d*")
	return func(value string) float64 {
		numericValue := re.FindString(value)
		unit := strings.Replace(value, numericValue, "", 1)
		convertedNum, _ := strconv.ParseFloat(numericValue, 64)
		multiplier, err := unitsMap[unit]
		if !err {
			panic("PARSED incorrectly")
		}
		return convertedNum * multiplier
	}
}
