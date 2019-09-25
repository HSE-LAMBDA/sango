package src

import (
	"encoding/json"
	lib "gosan"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func StartSimulation(flags *lib.SystemFlags) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lib.SIM_run(flags.SimRun)
}

func PlatformInit(flags *lib.SystemFlags) SHD {
	lib.SIM_init(flags.PacketFileName)
	parMap := SimParallelization(flags.NumJobsFileName)
	lib.SIM_platform_init(flags.PlatformFileName, flags.NDisks)
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
