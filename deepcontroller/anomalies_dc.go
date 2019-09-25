package deepcontroller

import (
	lib "gosan"
)

/*
========================
Common anomalies interface from DC
========================
*/

func (DCM *Manager) CreateAnomaly(P *lib.Process, params *ComponentAnomaly,
	component lib.DCAble, cTime float64) {
	if params == nil {
		return
	}

	var currentState = component.GetCurrentState()

	if currentState == params.Anomaly {
		return
	}
	if currentState == "default" {
		component.Break(P, cTime, params.Degree)
		component.SetCurrentState("nd")
		return
	}
	if currentState != "default" && params.Anomaly == "default" {
		component.Repair(P, cTime)
		component.SetCurrentState("default")
		return
	}
}
