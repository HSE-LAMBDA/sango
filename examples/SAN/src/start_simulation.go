package src

import (
	"encoding/json"
	lib "gosan"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func StartSimulation(flags *SystemFlags) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lib.SIM_run(flags.simRun)
}

func PlatformInit(flags *SystemFlags) SHD {
	lib.SIM_init(flags.packetFileName)
	parMap := SimParallelization(flags.numJobsFileName)
	lib.SIM_platform_init(flags.platformFileName, flags.nDisks)
	return init_SAN(flags, parMap)
}

func SimParallelization(filename string) map[int]float64 {
	parMap := make(map[int]float64)
	bytes := ParseFileAndUnmarshal(filename)

	err := json.Unmarshal(bytes, &parMap)
	if err != nil {
		panic(err)
	}

	return parMap

}

// go run main.go -platform= -deployment= -packet= -output= -file_size= -load_range= -anomaly_type= -anomaly_amount= -= -anomaly_time_range= -anomaly_duration=
