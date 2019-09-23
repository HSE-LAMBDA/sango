package src

import (
	"lib"
	"log"
)

type ComponentStatus int

const (
	FAIL ComponentStatus = iota
	OK

	// TODO Not implemented
	OKK
	MISSED
	BAD
	LOST
	MANUAL_RECOVERY
	RECOVERING
	RECOVERED
	FAIL_RECOVERY
	REPLACED
	DELETING
	DELETED
)

func (AM *AnomalyManager) AnomalyManagerProcess(TP *SANProcess, data interface{}) {
	args, ok := data.(*AnomalyFlags)
	if !ok {
		log.Panic("Anomaly flag data conversion error")
	}

	if args == nil || args.AnomalyAmount < 1 {
		return
	}
	kwargs := make(map[string]interface{})
	kwargs["process"] = TP

	AnomalyFunction, ok := AM.anomalies[args.AnomalyType]
	if !ok {
		log.Panic("No such anomaly. Stopping...")
	}

	WaitTime := fRand
	Duration := fRand

	for ; args.AnomalyAmount > 0; args.AnomalyAmount-- {
		TP.SIM_wait(WaitTime(args.MinAnomalyTime, args.MaxAnomalyTime))
		currentTime := lib.SIM_get_clock()
		kwargs["t1"] = currentTime
		kwargs["t2"] = currentTime + Duration(args.MinAnomalyDuration, args.MaxAnomalyDuration)
		AnomalyFunction(kwargs)
	}
}

type AnomalyManager struct {
	anomalies map[string]func(map[string]interface{})
}

func NewAnomalyManager() *AnomalyManager {
	am := &AnomalyManager{
		anomalies: make(map[string]func(map[string]interface{})),
	}
	am.anomalies["VESNIN1ANOMALY"] = Vesnin1Anomaly
	am.anomalies["VESNIN2ANOMALY"] = Vesnin2Anomaly
	am.anomalies["VESNIN1CLIENTLINK"] = Vesnin1ClientLinkDegradation
	am.anomalies["VESNIN2CLIENTLINK"] = Vesnin2ClientLinkDegradation

	return am
	//todo
	//anomalyCollection := make(map[string]map[string]*AnomalyType)

	//anomalyCollection["Controller"]["VESNIN1ANOMALY"] = Vesnin1Anomaly
	//anomalyCollection["Controller"]["VESNIN1ANOMALY"] = Vesnin2Anomaly
	//anomalyCollection["Link"]["VESNIN1CLIENTLINK"] = Vesnin1ClientLinkDegradation
	//anomalyCollection["Link"]["VESNIN2CLIENTLINK"] = Vesnin2ClientLinkDegradation
}

func Vesnin1Anomaly(kwargs map[string]interface{}) {
	kwargs["hostName"] = "Server1"
	ServerBreakAndRepair(kwargs)
}

func Vesnin2Anomaly(kwargs map[string]interface{}) {
	kwargs["hostName"] = "Server2"
	ServerBreakAndRepair(kwargs)
}

func Vesnin1ClientLinkDegradation(kwargs map[string]interface{}) {
	kwargs["linkName"] = "Client_Server1"
	kwargs["newState"] = 0.1
	LinkBreakAndRepair(kwargs)
}

func Vesnin2ClientLinkDegradation(kwargs map[string]interface{}) {
	kwargs["linkName"] = "Client_Server2"
	kwargs["newState"] = 0.1
	LinkBreakAndRepair(kwargs)
}

func ServerBreakAndRepair(kwargs map[string]interface{}) {
	process, ok := kwargs["process"].(*lib.Process)
	if !ok {
		log.Panic("No  process key anomalies")
	}
	hostName, ok := kwargs["hostName"].(string)
	if !ok {
		log.Panic("No hostname key anomalies")
	}
	host := lib.GetHostByName(hostName)
	t1, ok := kwargs["t1"].(float64)
	if !ok {
		log.Panic("No time start key anomalies")
	}
	t2, ok := kwargs["t2"].(float64)
	if !ok {
		log.Panic("No time end key anomalies")
	}
	process.CreateAnomalyHostSync(host, t1, 0)
	process.RepairHostSync(host, t2)
}

func LinkBreakAndRepair(kwargs map[string]interface{}) {
	process, ok := kwargs["process"].(*lib.Process)
	if !ok {
		log.Panic("No process key anomalies")
	}
	linkName, ok := kwargs["linkName"].(string)
	if !ok {
		log.Panic("No linkname key anomalies")
	}
	link := lib.GetLinkByName(linkName)
	t1, ok := kwargs["t1"].(float64)
	if !ok {
		log.Panic("No time start key anomalies")
	}
	t2, ok := kwargs["t2"].(float64)
	if !ok {
		log.Panic("No time end key anomalies")
	}
	newState, ok := kwargs["newState"].(float64)
	if !ok {
		log.Panic("No newState key anomalies")
	}

	process.CreateAnomalyLinkSync(link, t1, newState)
	process.RepairLinkSync(link, t2, 1)
}

type (
	AnomalyCollection map[string]map[string]*AnomalyType

	AnomalyType struct {
		Create func(float64, string, float64) `json:"create"`
		Repair func(float64, string)          `json:"repair"`
	}
)
