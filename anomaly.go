package lib

import "fmt"

var _ = fmt.Print

type Anomaly struct {
	resource Resource
	newState float64

	anomalyPart float64
}

func NewAnomaly(resource Resource, newState float64) (e *Anomaly) {
	e = &Anomaly{
		resource: resource,
		newState: newState,
	}
	return
}
