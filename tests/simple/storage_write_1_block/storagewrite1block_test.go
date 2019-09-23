package main

import (
	"fmt"
	"lib"
)

const KILO = 1000
const MEGA = 1000000

func ExampleWriteOneBlock() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("write_packet", write_1_packet)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Used space before: 0.00
	// Used space after: 2097.15
	// Simulation took 3.00
}

func write_1_packet(p *lib.Process, args []string) {
	storage := p.GetHost().GetOneStorage()
	fmt.Printf("Used space before: %.2f\n", storage.GetUsedSpace()*KILO)
	res := p.WriteSync(storage, lib.PACKET_2048K)
	if res != lib.OK {
		panic("Error")
	}
	fmt.Printf("Used space after: %.2f\n", storage.GetUsedSpace()*KILO)
}
