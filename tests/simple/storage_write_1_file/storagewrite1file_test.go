package main

import (
	"fmt"
	"lib"
)

const KILO = 1000
const MEGA = 1000000

func ExampleWriteOneFile() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("write_packet", write_1_packet)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Used space before: 0.00 KB
	// Used space after: 102.40 KB
	// Simulation took 75.00
}

func write_1_packet(p *lib.Process, args []string) {
	file := lib.NewFile("1", 100*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size

	storage := p.GetHost().GetOneStorage()
	fmt.Printf("Used space before: %.2f KB\n", storage.GetUsedSpace()*KILO)

	for n > 0 {
		res := p.WriteSync(storage, lib.PACKET_4K)
		if res != lib.OK {
			panic("Error")
		}
		n--
	}

	fmt.Printf("Used space after: %.2f KB\n", storage.GetUsedSpace()*KILO)
}
