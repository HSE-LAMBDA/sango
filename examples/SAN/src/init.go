package src

import (
	"encoding/json"
	lib "gosan"
	"gosan/deepcontroller"
	"strings"
)

func init_SAN(flags *lib.SystemFlags, parMap map[int]float64) SHD {
	controllers := initControllers()
	jbods := initSANJbods()

	iob := NewIOBalancer(lib.GetHostByName("Helper"), controllers, jbods)
	cm := initClient(iob, flags, parMap)
	am := initAnomalyManager(flags.AnomalyFlags)
	acm := initAtmosphereControlManager(flags.AtmControlFileName)
	afm := initAtmosphereFailManager(flags.AtmDependenciesFileName, acm, iob)
	vols := initVolumes(iob)

	tm := initTracer(flags.TracerFlags, iob, acm)
	dcm := deepcontroller.Init_d(flags.DCFlags)

	initComponentsWithIOBalancer(iob, tm.logs.StState, controllers)

	assembly := NewSAN()
	san := assembly.IOBalancer(iob).Controllers(controllers).JbodControllers(jbods).Tracer(tm).AtmosphereControl(acm).
		Build()

	_ = am
	_ = afm
	_ = dcm
	_ = vols
	_ = cm
	return san

}


func NewIOBalancer(host *lib.Host, cons []*Controller, jbods []*SANJBODController) *IOBalancer {
		naming := NewNamingProps(host.Name, host.Type, host.Id)
		iob := &IOBalancer{
		NamingProps: naming,
		Host:        host,
		IOBalancerProps: &IOBalancerProps{
		CommonProps: &CommonProps{
		Status: OK,
	},
	},
		Controllers: cons,
		Jbods:jbods,
		controllersMap: make(map[string]*Controller),
		openedFiles:   make(map[string]*FileInfo),
	}

		return iob
}


func initClient(iob *IOBalancer, flags *lib.SystemFlags, parMap map[int]float64) map[int]*ClientManager {
	clientHost := lib.GetHostByName("Client")
	clients := make(map[int]*ClientManager)

	filename := flags.ClientFileName
	writerFlags := flags.WriteClientFlags
	readerFlags := flags.ReadClientFlags

	if filename != "" {

		clientMap := make(map[string]*ClientManager)
		bytes := ParseFileAndUnmarshal(filename)
		err := json.Unmarshal(bytes, &clientMap)

		if err != nil {
			panic(err)
		}

		for id := range clientMap {

			cm := NewClientManager(clientHost, iob, id, 1e7, parMap)
			_ = cm
			if id == "dc" {
				//dcClient = cm
				continue
			}
			//FORK("writerClientManager", cm.ClientManagerAsyncProcess, clientHost, writerFlags)
			//FORK("readerClientManager", cm.ClientManagerAsyncProcess, clientHost, readerFlags)

			// FORK("clientTracer", cm.Tracer, clientHost, readerFlags)
		}
	} else {
		n := flags.NumJobs
		for n > 0 {
			c := NewClientManager(clientHost, iob, "dc-client", 1e70, parMap)
			clients[n] = c
			n--
		}
	}

	master := NewMaster(clientHost, clients)
	FORK("writerClientManager", master.ClientManagerAsyncProcess, clientHost, writerFlags)
	FORK("readerClientManager", master.ClientManagerAsyncProcess, clientHost, readerFlags)

	return clients
}

func initControllers() []*Controller {
	hosts := lib.GetHosts()
	cons := make([]*Controller, 0)

	index := 0
	for hostName, host := range hosts {
		if strings.HasPrefix(hostName, "Server") {
			con := NewController(host)
			con.index = index
			cons = append(cons, con)
			index++
		}
	}
	return cons
}

func initSANJbods() []*SANJBODController {
	jbods := lib.GetAllJBODs()
	tatJBODs := make([]*SANJBODController, 0)
	i := 0
	for name, libJBOD := range jbods {
		if strings.HasPrefix(name, "JBOD") {
			jbod := NewSANJBODController(libJBOD)
			tatJBODs = append(tatJBODs, jbod)
			i++
		}
	}

	return tatJBODs
}


func initTracer(data *lib.TracerFlags, iob *IOBalancer, acm *AtmosphereControlManager) *TracerManager {
	tm := NewTracerManager(iob, acm, nil)
	host := lib.GetHostByName("Helper")
	FORK("TracerM", tm.TracerManagerProcess, host, data)
	return tm
}


func initAnomalyManager(data *lib.AnomalyFlags) *AnomalyManager {
	am := NewAnomalyManager()
	host := lib.GetHostByName("Helper")
	FORK("AM", am.AnomalyManagerProcess, host, data)
	return am
}

func initAtmosphereControlManager(fileName string) *AtmosphereControlManager {
	acm := NewAtmosphereControlManager(fileName)
	host := lib.GetHostByName("Helper")
	FORK("ACM", acm.AtmosphereControlManagerProcess, host, fileName)
	return acm
}

func initAtmosphereFailManager(fileName string, acm *AtmosphereControlManager, iob *IOBalancer) *AtmosphereFailManager {
	afm := NewAtmosphereFailManager(acm, iob)
	iob.FailManager = afm

	host := lib.GetHostByName("Helper")
	FORK("BreakDisks", afm.AtmosphereDiskFailManagerProcess, host, fileName)

	conMap := make(map[string]lib.DCAble)
	linkMap := make(map[string]lib.DCAble)

	for name, con := range iob.controllersMap {
		conMap[name] = con
	}

	for name, link := range iob.linksMap {
		linkMap[name] = link
	}

	FORK("BreakCons", afm.PoissonBreakManagerProcess, host, conMap)
	FORK("BreakLinks", afm.PoissonBreakManagerProcess, host, linkMap)

	FORK("Repair", afm.PoissonRepairManagerProcess, host, nil)
	return afm
}

func initComponentsWithIOBalancer(iob *IOBalancer, gds *DiskState, cons []*Controller) {
	// cons := iob.Controllers
	// cacheCons := iob.Caches
	jbodCons := iob.Jbods

	for _, jbodCon := range jbodCons {
		disks := jbodCon.disks
		for _, disk := range disks {
			disk.iob = iob
			disk.GlobalDiskState = gds
		}
	}
	iob.GlobalDiskState = gds

	for _, con := range cons {
		con.iob = iob
	}
}
