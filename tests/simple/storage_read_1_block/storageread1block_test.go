package main

import (
	"fmt"
	"lib"
)

const KILO = 1000
const MEGA = 1000000

func ExampleReadOneBlock() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("read_packet", read_1_packet)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Operation took: 3.00
	// Simulation took 3.00
}

func read_1_packet(p *lib.Process, args []string) {
	host := p.GetHost()
	storage := host.GetOneStorage()
	p.Read(storage, lib.PACKET_4K)
	fmt.Printf("Operation took: %.2f\n", lib.SIM_get_clock())
}
