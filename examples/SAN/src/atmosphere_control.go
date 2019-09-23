package src

import (
	"encoding/json"
	"lib"
	"log"
	"sort"
	"strconv"
)

type (
	AtmosphereControl struct {
		Temp      float64 `json:"air_temp"`
		Humidity  float64 `json:"humidity"`
		Pressure  float64 `json:"atm_pressure"`
		Vibration float64 `json:"vibration"`
	}
	AtmosphereControlManager struct {
		currentAtmosphere *AtmosphereControl
		atmTimeControlMap map[float64]*AtmosphereControl
	}
)

func (ACM *AtmosphereControlManager) AtmosphereControlManagerProcess(TP *SANProcess, data interface{}) {
	atmControlFileName, ok := data.(string)
	if !ok {
		log.Panic("atmControlFileName error")
	}

	if atmControlFileName == "" || ACM.atmTimeControlMap == nil {
		log.Println("No atmosphere control management mode")
		return
	}

	previousTime := lib.SIM_get_clock()

	var times []float64
	for key := range ACM.atmTimeControlMap {
		times = append(times, key)
	}
	sort.Float64s(times)

	for _, timeX := range times {
		currentAtm, ok := ACM.atmTimeControlMap[timeX]
		if !ok {
			log.Panic("Atmosphere problem")
		}
		TP.SIM_wait(timeX - previousTime)
		*ACM.currentAtmosphere = *currentAtm
		previousTime = timeX
	}
}

func NewAtmosphereControlManager(atmControlFileName string) *AtmosphereControlManager {
	var atmTimeControlMap map[float64]*AtmosphereControl
	if atmControlFileName != "" {
		atmTimeControlMap = ParseAtmosphereControlFile(atmControlFileName)
	}

	return &AtmosphereControlManager{
		currentAtmosphere: &AtmosphereControl{},
		atmTimeControlMap: atmTimeControlMap,
	}
}

func ParseAtmosphereControlFile(filename string) map[float64]*AtmosphereControl {
	sMap := make(map[string]*AtmosphereControl)
	bytes := ParseFileAndUnmarshal(filename)

	err := json.Unmarshal(bytes, &sMap)

	if err != nil {
		panic(err)
	}

	return ConvertStringToFloat64(sMap)
}

func ConvertStringToFloat64(sMap map[string]*AtmosphereControl) map[float64]*AtmosphereControl {
	fMap := make(map[float64]*AtmosphereControl)
	for name, value := range sMap {
		f64, err := strconv.ParseFloat(name, 64)
		if err != nil {
			panic("conversion from string to float64 error")
		}
		fMap[f64] = value
	}
	sMap = nil
	return fMap
}
