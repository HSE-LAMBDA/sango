package src

import (
	"fmt"
	lib "gosan"
	"log"
	"math"
	"math/rand"
)

type AtmosphereFailManager struct {
	currentAtmosphere *AtmosphereControl
	iob               *IOBalancer
	br                []lib.DCAble
	linkMap, conMap   map[string]lib.DCAble
}

func (AFM *AtmosphereFailManager) AtmosphereDiskFailManagerProcess(TP *SANProcess, data interface{}) {
	atmDependenciesFileName, ok := data.(string)
	if !ok {
		log.Panic("atmDependenciesFileName error")
	}

	if atmDependenciesFileName == "" {
		log.Println("No atmosphere dependencies file")
		return
	}

	TP.Daemonize()

	atmDepMap := ParseAtmosphereInputDependency(atmDependenciesFileName)

	tempRoot := GetAtmosphereDependencyRoot("Temperature", atmDepMap)
	humidityRoot := GetAtmosphereDependencyRoot("Humidity", atmDepMap)
	pressureRoot := GetAtmosphereDependencyRoot("Pressure", atmDepMap)
	vibrationRoot := GetAtmosphereDependencyRoot("Vibration", atmDepMap)

	disks := AFM.iob.disksMap

	for {
		for _, disk := range disks {
			atmControl := AFM.currentAtmosphere
			f1 := tempRoot.getFailProbability(atmControl.Temp)
			f2 := humidityRoot.getFailProbability(atmControl.Humidity)
			f3 := pressureRoot.getFailProbability(atmControl.Pressure)
			f4 := vibrationRoot.getFailProbability(atmControl.Vibration)

			failure := f1 * (f2 + f3 + f4)

			if ok := MonteCarlo(failure); ok {
				TP.CreateDiskAnomalyFullOFF(lib.SIM_get_clock(), disk.Storage, 1)
				AFM.br = append(AFM.br, disk)
			}
		}
		TP.SIM_wait(0.1)
	}
}

func (AFM *AtmosphereFailManager) PoissonBreakManagerProcess(TP *SANProcess, data interface{}) {
	return
	//TP.Daemonize()
	//components, ok := data.(map[string]lib.DCAble)
	//if !ok {
	//	log.Panic("No such error")
	//}
	//lambda := 42.
	//seed := int64(42)
	//poisson := NewPoissonGenerator(seed)
	//
	//for {
	//	for _, component := range components {
	//		component.Break(TP, lib.SIM_get_clock(), 1)
	//		AFM.br = append(AFM.br, component)
	//	}
	//	TP.SIM_wait(float64(poisson.Poisson(lambda)))
	//}
}

func (AFM *AtmosphereFailManager) PoissonRepairManagerProcess(TP *SANProcess, data interface{}) {
	return
	//TP.Daemonize()
	//
	//lambda := 42.
	//seed := int64(42)
	//poisson := NewPoissonGenerator(seed)
	//
	//for {
	//	for _, component := range AFM.br {
	//		component.Repair(TP, 1)
	//	}
	//	TP.SIM_wait(float64(poisson.Poisson(lambda)))
	//}
}

func NewAtmosphereFailManager(acm *AtmosphereControlManager, iob *IOBalancer) *AtmosphereFailManager {
	return &AtmosphereFailManager{
		currentAtmosphere: acm.currentAtmosphere,
		iob:               iob,
	}
}

func (ad *AtmosphereDependency) getFailProbability(current float64) float64 {
	closestValue := binarySearch(ad.arr, current)
	value, ok := ad.table[closestValue]
	if !ok {
		panic("No such value in table")
	}
	return value
}

type UniformGenerator struct {
	rd *rand.Rand
}

func NewUniformGenerator(seed int64) *UniformGenerator {
	return &UniformGenerator{rd: rand.New(rand.NewSource(seed))}
}

type PoissonGenerator struct {
	uniform *UniformGenerator
}

func NewPoissonGenerator(seed int64) *PoissonGenerator {
	urng := NewUniformGenerator(seed)
	return &PoissonGenerator{urng}
}

// Poisson returns a random number of possion distribution
func (prng *PoissonGenerator) Poisson(lambda float64) int64 {
	if !(lambda > 0.0) {
		panic(fmt.Sprintf("Invalid lambda: %.2f", lambda))
	}
	return prng.poisson(lambda)
}

func (prng *PoissonGenerator) poisson(lambda float64) int64 {
	// algorithm given by Knuth
	L := math.Pow(math.E, -lambda)
	var k int64 = 0
	var p float64 = 1.0

	for p > L {
		k++
		p *= prng.uniform.rd.Float64()
	}
	return k - 1
}
