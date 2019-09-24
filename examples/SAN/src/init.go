package src

import (
	"encoding/json"
	lib "gosan"
	"strings"
)

func init_SAN(flags *SystemFlags, parMap map[int]float64) SHD {
	controllers := initControllers()
	jbods := initSANJbods()

	iob := initIobalancer(controllers, jbods)
	cm := initClient(iob, flags, parMap)
	am := initAnomalyManager(flags.anomalyFlags)
	acm := initAtmosphereControlManager(flags.atmControlFileName)
	afm := initAtmosphereFailManager(flags.atmDependenciesFileName, acm, iob)
	vols := initVolumes(iob)

	tm := initTracer(flags.tracerFlags, iob, acm)
	dcm := initDeepController(flags.dcFlags, iob, tm, cm)

	initComponentsWithIOBalancer(iob, tm.logs.StState, controllers)

	assembly := NewSAN()
	san := assembly.IOBalancer(iob).Controllers(controllers).JbodControllers(jbods).Tracer(tm).AtmosphereControl(acm).
		Build()

	_ = am
	_ = afm
	_ = dcm
	_ = vols
	return san

}

func initClient(iob *IOBalancer, flags *SystemFlags, parMap map[int]float64) map[int]*ClientManager {
	clientHost := lib.GetHostByName("Client")
	clients := make(map[int]*ClientManager)

	filename := flags.clientFileName
	writerFlags := flags.writeClientFlags
	readerFlags := flags.readClientFlags

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
		n := flags.numJobs
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

func initIobalancer(cons []*Controller, jbods []*SANJBODController) *IOBalancer {
	host := lib.GetHostByName("IOBalancer")
	iob := NewIOBalancer(host, cons, jbods)
	FORK("iob", iob.IOBalancerProcessManager, host, nil)
	return iob
}

func initTracer(data *TracerFlags, iob *IOBalancer, acm *AtmosphereControlManager) *TracerManager {
	tm := NewTracerManager(iob, acm, nil)
	host := lib.GetHostByName("Helper")
	FORK("TracerM", tm.TracerManagerProcess, host, data)
	return tm
}

func initDeepController(data *DeepControllerFlags, iob *IOBalancer, t *TracerManager,
	cm map[int]*ClientManager) *DeepControllerManager {
	dcm := NewDeepControllerManager(data, iob, cm, t)
	host := lib.GetHostByName("Helper")
	FORK("DCM", dcm.DeepControllerManagerProcess, host, data)

	rwConversion = rwConversionFactory()
	return dcm
}

func initAnomalyManager(data *AnomalyFlags) *AnomalyManager {
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

	conMap := make(map[string]BreakAble)
	linkMap := make(map[string]BreakAble)

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
