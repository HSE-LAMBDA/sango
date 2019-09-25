package lib

import (
	"flag"
	"strings"
)

type (
	SystemFlags struct {
		PlatformFileName        string
		PacketFileName          string
		AtmControlFileName      string
		AtmDependenciesFileName string
		ClientFileName          string
		NumJobsFileName         string

		SimRun  float64
		NDisks  int
		NumJobs int

		WriteClientFlags *ClientFlags
		ReadClientFlags  *ClientFlags
		AnomalyFlags     *AnomalyFlags
		TracerFlags      *TracerFlags
		DCFlags          *DeepControllerFlags
	}

	ClientFlags struct {
		FileAmount   int
		MinFileSize  float64
		MaxFileSize  float64
		MinPauseTime float64
		MaxPauseTime float64
		RequestType  RequestType

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


	TracerFlags struct {
		OutputFileName string
	}

	 DeepControllerFlags struct {
		DeepControllingMode bool
		DeepHost            string
		DeepPort            string
		DeepProtocol        string
		DeepJsonDelay       float64

		SimRun float64
	}
)

func NewSystemFlags() *SystemFlags {
	sf := &SystemFlags{
		WriteClientFlags: &ClientFlags{
			RequestType: WRITE,
		},
		ReadClientFlags: &ClientFlags{
			RequestType: READ,
		},
		AnomalyFlags: new(AnomalyFlags),
		TracerFlags:  new(TracerFlags),
		DCFlags:      new(DeepControllerFlags),
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
	flag.StringVar(&sf.PlatformFileName, "platform", "", "topology filename")
	flag.StringVar(&sf.PacketFileName, "packet", "", "the packet configs")
	flag.StringVar(&sf.AtmControlFileName, "atm_control", "", "the atm control configs")
	flag.StringVar(&sf.AtmDependenciesFileName, "atm_dep", "", "the atm dependencies configs")
	flag.StringVar(&sf.TracerFlags.OutputFileName, "output", "", "the filename")
	flag.StringVar(&sf.ClientFileName, "client", "", "client config")
	flag.StringVar(&sf.NumJobsFileName, "num_jobs_config", "", "num jobs config")

	// write mode
	flag.IntVar(&sf.WriteClientFlags.FileAmount, "file_amount_w", 0, "write fileAmount")
	flag.StringVar(&sf.WriteClientFlags.fileSizeRangeString, "file_size_w", "1MB..1MB", "write fileSizeRange")
	flag.StringVar(&sf.WriteClientFlags.loadRangeString, "load_range_w", "0s..0s", "write loadRange")

	// read mode
	flag.IntVar(&sf.ReadClientFlags.FileAmount, "file_amount_r", 0, "read fileAmount")
	flag.StringVar(&sf.ReadClientFlags.fileSizeRangeString, "file_size_r", "1MB..1MB", "read fileSizeRange")
	flag.StringVar(&sf.ReadClientFlags.loadRangeString, "load_range_r", "0s..0s", "read loadRange")

	flag.StringVar(&sf.AnomalyFlags.AnomalyType, "anomaly_type", "VESNIN1ANOMALY", "anomalyType")
	flag.IntVar(&sf.AnomalyFlags.AnomalyAmount, "anomaly_amount", 0, "anomalyAmount")
	flag.StringVar(&sf.AnomalyFlags.anomalyTimeRangeString, "anomaly_time_range", "20s..30s", "anomalyTimeRange")
	flag.StringVar(&sf.AnomalyFlags.anomalyDurationString, "anomaly_duration", "1s..5s", "anomalyDuration")

	// additional values to ~ deep controller
	flag.BoolVar(&sf.DCFlags.DeepControllingMode, "controlling_mode", false, "controllingMode")
	flag.StringVar(&sf.DCFlags.DeepHost, "host", "localhost", "host")
	flag.StringVar(&sf.DCFlags.DeepPort, "port", "1337", "port")
	flag.StringVar(&sf.DCFlags.DeepProtocol, "protocol", "tcp", "tcp")
	flag.Float64Var(&sf.DCFlags.DeepJsonDelay, "json_delay", 1, "json delay")

	flag.Float64Var(&sf.SimRun, "sim_run", -1, "end of simulation")
	flag.IntVar(&sf.NDisks, "disk_amount", -1, "number of disks")
	flag.IntVar(&sf.NumJobs, "num_jobs", -1, "number of jobs")

	flag.Parse()

	sf.DCFlags.SimRun = sf.SimRun
	return sf
}

func ParseFlags(sf *SystemFlags) *SystemFlags {
	sf.WriteClientFlags.MinFileSize, sf.WriteClientFlags.MaxFileSize = RangeParser(sf.WriteClientFlags.
		fileSizeRangeString)
	sf.WriteClientFlags.MinPauseTime, sf.WriteClientFlags.MaxPauseTime = RangeParser(sf.WriteClientFlags.
		loadRangeString)

	sf.ReadClientFlags.MinFileSize, sf.ReadClientFlags.MaxFileSize = RangeParser(sf.ReadClientFlags.fileSizeRangeString)
	sf.ReadClientFlags.MinPauseTime, sf.ReadClientFlags.MaxPauseTime = RangeParser(sf.ReadClientFlags.loadRangeString)

	sf.AnomalyFlags.MinAnomalyTime, sf.AnomalyFlags.MaxAnomalyTime = RangeParser(sf.AnomalyFlags.anomalyTimeRangeString)
	sf.AnomalyFlags.MinAnomalyDuration, sf.AnomalyFlags.MaxAnomalyDuration = RangeParser(sf.AnomalyFlags.
		anomalyDurationString)

	return sf
}

func RangeParser(flag string) (x1, x2 float64) {
	slice := strings.Split(flag, "..")
	x1 = UnitToFloatParser(slice[0])
	x2 = UnitToFloatParser(slice[1])
	return
}
