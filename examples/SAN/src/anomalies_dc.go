package src

import (
	lib "gosan"
)

/*
========================
Common anomalies interface from DC
========================
*/

type BreakAble interface {
	Break(*SANProcess, float64, float64)
	Repair(*SANProcess, float64)
	GetCurrentState() string
	SetCurrentState(string)
}

func (DCM *DeepControllerManager) CreateAnomaly(TP *SANProcess, params *ComponentAnomaly,
	component BreakAble, cTime float64) {
	if params == nil {
		return
	}

	var currentState = component.GetCurrentState()

	if currentState == params.Anomaly {
		return
	}
	if currentState == "default" {
		component.Break(TP, cTime, params.Degree)
		component.SetCurrentState("nd")
		return
	}
	if currentState != "default" && params.Anomaly == "default" {
		component.Repair(TP, cTime)
		component.SetCurrentState("default")
		return
	}
}

/*
========================
SAN Component anomaly
========================
*/

func (tc *SANComponent) Break(TP *SANProcess, currentTime, degree float64) {}

func (tc *SANComponent) Repair(TP *SANProcess, currentTime float64) {}

/*
========================
Controller anomaly
========================
*/

func (controller *Controller) Break(TP *SANProcess, currentTime, degree float64) {
	// delete controller from array of controllers
	controller.iob.RemoveWorkingController(controller.index)

	lib.CreateAnomalyHostAsync(TP.Process, controller.Host, currentTime, 1-degree)
}

func (controller *Controller) Repair(TP *SANProcess, currentTime float64) {
	iob := controller.iob

	// add controller to the list of controllers
	iob.AddWorkingController(controller)

	lib.RepairHostAsync(TP.Process, controller.Host, currentTime)
	controller.currentState = "default"

	iob.FailManager.conMap[controller.ID] = controller
}

/*
========================
Disk anomaly
========================
*/

func (disk *SANDisk) Break(TP *SANProcess, currentTime, degree float64) {

	// delete disk from array of disks
	disk.jbod.disksSlice[disk.index] = disk.jbod.disksSlice[len(disk.jbod.disksSlice)-1]
	disk.jbod.disksSlice[len(disk.jbod.disksSlice)-1] = nil
	disk.jbod.disksSlice = disk.jbod.disksSlice[:len(disk.jbod.disksSlice)-1]

	lib.CreateDiskAnomalyAsync(TP.Process, disk.Storage, currentTime+0.005, 1-degree)
}

func (disk *SANDisk) Repair(TP *SANProcess, currentTime float64) {

	// delete disk from array of disks
	index := len(disk.jbod.disksSlice)
	disk.index = index
	disk.jbod.disksSlice = append(disk.jbod.disksSlice, disk)
}

/*
========================
Link anomaly
========================
*/

func (link *SANLink) Break(TP *SANProcess, currentTime, degree float64) {

	for _, node := range link.srcDst {
		switch node.(type) {
		case *Controller:
			con, _ := node.(*Controller)
			con.iob.RemoveWorkingController(con.index)
		default:
			_ = 42
		}
	}

	lib.CreateAnomalyLinkAsync(TP.Process, link.Link, currentTime, 1-degree)
}

func (link *SANLink) Repair(TP *SANProcess, currentTime float64) {
	lib.RepairLinkAsync(TP.Process, link.Link, currentTime)

	link.iob.FailManager.linkMap[link.ID] = link

	// add controller after link get repaired
	for _, node := range link.srcDst {
		switch node.(type) {
		case *Controller:
			con, _ := node.(*Controller)
			con.iob.AddWorkingController(con)
		default:
			_ = 42
		}
	}
}

/*
========================
JBOD anomaly
========================
*/

func (jbod *SANJBODController) Break(TP *SANProcess, currentTime, degree float64) {}

func (jbod *SANJBODController) Repair(TP *SANProcess, currentTime float64) {}
