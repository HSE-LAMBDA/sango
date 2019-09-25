package anomaly

import (
	"fmt"
	lib "gosan"
)

var _ = fmt.Print

type Anomaly struct {
	resource lib.Resource
	newState float64

	anomalyPart float64
}

func NewAnomaly(resource lib.Resource, newState float64) (e *Anomaly) {
	e = &Anomaly{
		resource: resource,
		newState: newState,
	}
	return
}
