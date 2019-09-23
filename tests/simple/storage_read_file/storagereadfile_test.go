package main

import (
	"lib"
)

const KILO = 1000
const MEGA = 1000000

func ExampleReadOneFile() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("read_packet", read_1_packet)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Simulation took 75.00
}

func read_1_packet(p *lib.Process, args []string) {
	file := lib.NewFile("1", 100*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size

	host := p.GetHost()
	storage := host.GetOneStorage()

	for n > 0 {
		p.Read(storage, lib.PACKET_4K)
		n--
	}
}
