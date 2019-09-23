package src

import (
	"flag"
	"lib"
	"strings"
)

type (
	SystemFlags struct {
		platformFileName        string
		packetFileName          string
		atmControlFileName      string
		atmDependenciesFileName string
		clientFileName          string
		numJobsFileName         string

		simRun  float64
		nDisks  int
		numJobs int

		writeClientFlags *ClientFlags
		readClientFlags  *ClientFlags
		anomalyFlags     *AnomalyFlags
		tracerFlags      *TracerFlags
		dcFlags          *DeepControllerFlags
	}

	ClientFlags struct {
		FileAmount   int
		MinFileSize  float64
		MaxFileSize  float64
		MinPauseTime float64
		MaxPauseTime float64
		RequestType  lib.RequestType

		fileSizeRangeString string
		loadRangeString     string
	}

	AnomalyFlags struct {
		AnomalyType        string
		AnomalyAmount      int
		MinAnomalyTime     float64
		MaxAnomalyTime     float64
		MinAnomalyDuration float64
		MaxAnomalyDuration float64

		anomalyTimeRangeString string
		anomalyDurationString  string
	}

	DeepControllerFlags struct {
		DeepControllingMode bool
		DeepHost            string
		DeepPort            string
		DeepProtocol        string
		DeepJsonDelay       float64

		SimRun float64
	}

	TracerFlags struct {
		outputFileName string
	}
)

func NewSystemFlags() *SystemFlags {
	sf := &SystemFlags{
		writeClientFlags: &ClientFlags{
			RequestType: lib.WRITE,
		},
		readClientFlags: &ClientFlags{
			RequestType: lib.READ,
		},
		anomalyFlags: new(AnomalyFlags),
		tracerFlags:  new(TracerFlags),
		dcFlags:      new(DeepControllerFlags),
	}
	return sf
}

// platformFileName —  topology/platform.xml
// deploymentFileName —  topology/deployment.xml
// fileSizeRange —  5MB..100MB
// loadRange—  1s..2s
// anomalyType — VESNIN1FAIL
// anomalyNumber — 5
// anomalyTimeRange — 20s..30s
// anomalyTimeDuration — 1s..5s

func InitFlags(sf *SystemFlags) *SystemFlags {
	flag.StringVar(&sf.platformFileName, "platform", "", "topology filename")
	flag.StringVar(&sf.packetFileName, "packet", "", "the packet configs")
	flag.StringVar(&sf.atmControlFileName, "atm_control", "", "the atm control configs")
	flag.StringVar(&sf.atmDependenciesFileName, "atm_dep", "", "the atm dependencies configs")
	flag.StringVar(&sf.tracerFlags.outputFileName, "output", "", "the filename")
	flag.StringVar(&sf.clientFileName, "client", "", "client config")
	flag.StringVar(&sf.numJobsFileName, "num_jobs_config", "", "num jobs config")

	// write mode
	flag.IntVar(&sf.writeClientFlags.FileAmount, "file_amount_w", 0, "write fileAmount")
	flag.StringVar(&sf.writeClientFlags.fileSizeRangeString, "file_size_w", "1MB..1MB", "write fileSizeRange")
	flag.StringVar(&sf.writeClientFlags.loadRangeString, "load_range_w", "0s..0s", "write loadRange")

	// read mode
	flag.IntVar(&sf.readClientFlags.FileAmount, "file_amount_r", 0, "read fileAmount")
	flag.StringVar(&sf.readClientFlags.fileSizeRangeString, "file_size_r", "1MB..1MB", "read fileSizeRange")
	flag.StringVar(&sf.readClientFlags.loadRangeString, "load_range_r", "0s..0s", "read loadRange")

	flag.StringVar(&sf.anomalyFlags.AnomalyType, "anomaly_type", "VESNIN1ANOMALY", "anomalyType")
	flag.IntVar(&sf.anomalyFlags.AnomalyAmount, "anomaly_amount", 0, "anomalyAmount")
	flag.StringVar(&sf.anomalyFlags.anomalyTimeRangeString, "anomaly_time_range", "20s..30s", "anomalyTimeRange")
	flag.StringVar(&sf.anomalyFlags.anomalyDurationString, "anomaly_duration", "1s..5s", "anomalyDuration")

	// additional values to ~ deep controller
	flag.BoolVar(&sf.dcFlags.DeepControllingMode, "controlling_mode", false, "controllingMode")
	flag.StringVar(&sf.dcFlags.DeepHost, "host", "localhost", "host")
	flag.StringVar(&sf.dcFlags.DeepPort, "port", "1337", "port")
	flag.StringVar(&sf.dcFlags.DeepProtocol, "protocol", "tcp", "tcp")
	flag.Float64Var(&sf.dcFlags.DeepJsonDelay, "json_delay", 1, "json delay")

	flag.Float64Var(&sf.simRun, "sim_run", -1, "end of simulation")
	flag.IntVar(&sf.nDisks, "disk_amount", -1, "number of disks")
	flag.IntVar(&sf.numJobs, "num_jobs", -1, "number of jobs")

	flag.Parse()

	sf.dcFlags.SimRun = sf.simRun
	return sf
}

func ParseFlags(sf *SystemFlags) *SystemFlags {
	sf.writeClientFlags.MinFileSize, sf.writeClientFlags.MaxFileSize = RangeParser(sf.writeClientFlags.fileSizeRangeString)
	sf.writeClientFlags.MinPauseTime, sf.writeClientFlags.MaxPauseTime = RangeParser(sf.writeClientFlags.loadRangeString)

	sf.readClientFlags.MinFileSize, sf.readClientFlags.MaxFileSize = RangeParser(sf.readClientFlags.fileSizeRangeString)
	sf.readClientFlags.MinPauseTime, sf.readClientFlags.MaxPauseTime = RangeParser(sf.readClientFlags.loadRangeString)

	sf.anomalyFlags.MinAnomalyTime, sf.anomalyFlags.MaxAnomalyTime = RangeParser(sf.anomalyFlags.anomalyTimeRangeString)
	sf.anomalyFlags.MinAnomalyDuration, sf.anomalyFlags.MaxAnomalyDuration = RangeParser(sf.anomalyFlags.anomalyDurationString)

	return sf
}

func RangeParser(flag string) (x1, x2 float64) {
	slice := strings.Split(flag, "..")
	x1 = lib.UnitToFloatParser(slice[0])
	x2 = lib.UnitToFloatParser(slice[1])
	return
}
