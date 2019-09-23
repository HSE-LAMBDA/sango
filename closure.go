package lib

/*
============================
Closures
============================
*/

var SIM_get_clock func() float64
var SIM_run func(float64)
var GetCallbacks func(EventType) []func(*Event)
var NewProcess func(string, *Host) *Process

var PIDNext = PIDNextFactory()
var CIDNext = PIDNextFactory()
var UnitToFloatParser = UnitToFloatFactory()

/*
============================
Links
============================
*/

var GetLinkByName func(name string) *Link
var GetAllLinksMap func() map[string]*Link
var GetLinkBetweenHosts func(*Host, *Host) *Link

func initLinkMapClosures(linksMap map[string]*Link, routesMap, bRoutesMap map[*Host]map[*Host]*Link) {
	GetLinkByName = getLinkByNameFactory(linksMap)
	GetAllLinksMap = getAllLinksMapFactory(linksMap)
	GetLinkBetweenHosts = getLinkBetweenHostsFactory(routesMap, bRoutesMap)
}

/*
===================================
JBODs
===================================
*/

var GetAllJBODs func() map[string]*JBOD
var GetAllJBODsSlice func() []*JBOD

func initJBODClosures(jbodMap map[string]*JBOD, jSlice []*JBOD) {
	GetAllJBODs = getAllJBODsFactory(jbodMap)
	GetAllJBODsSlice = getAllJBODsSliceFactory(jSlice)
}

var GetCacheJbod func() *JBOD

func initGetCacheJbodClosure(cacheJbod *JBOD) {
	GetCacheJbod = getCacheJbodFactory(cacheJbod)
}

/*
============================
Volume functions
============================
*/
var GetAllVolumes func() map[string]*Volume

func initGetVolumeClosure(volumes []*Volume) {
	GetAllVolumes = getAllVolumesFactory(volumes)
}

/*
============================
Hosts
============================
*/

var GetHostByName func(string) *Host
var GetHosts func() map[string]*Host
var GetLinks func(*Host) []*Link

func initHostClosures(hostsMap map[string]*Host) {
	GetHostByName = getHostByNameFactory(hostsMap)
	GetHosts = getHostsFactory(hostsMap)
}

func initHostLinkClosures(hostsLinksMap map[*Host][]*Link) {
	GetLinks = getLinksFactory(hostsLinksMap)
}

/*
============================
Packets
============================
*/

var GetRandomPacket func(RequestType) *Packet
var GetRandomBlockSize func() string
var GetPacketByName func(RequestType, string) *Packet
var _basic_event_adding func(*Process, Resource, *Packet, EventType) (STATUS, float64)

func initPacketClosures(packetInfo map[RequestType]map[string]*Packet, packetIndices []string) {
	GetRandomPacket, GetRandomBlockSize = getRandomPacketFactory(packetInfo, packetIndices)
	GetPacketByName = getPacketByNameFactory(packetInfo)

	_basic_event_adding = _basic_event_adding_factory()
}

/*
============================
Deployment functions
============================
*/

var GetFunctionByName func(string) BaseBehaviourFunction
var SIM_function_register func(string, BaseBehaviourFunction, interface{})

func initDeploymentClosure(funcMap map[string]BaseBehaviourFunction, funcData map[string]interface{}) {
	SIM_function_register = sim_function_registerFactory(funcMap, funcData)
	GetFunctionByName = getFunctionByNameFactory(funcMap)
	_ = funcData
}
