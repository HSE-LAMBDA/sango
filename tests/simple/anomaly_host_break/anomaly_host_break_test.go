package main

import (
	"fmt"
	"lib"
)

func ExampleAnomalyHostBreak() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)
	lib.SIM_function_register("anomaly", anomaly)

	lib.SIM_launch_application("topology/deployment.xml")

	lib.SIM_run(nil)
	// Output:
	// FAIL sender. Quit at  7
	// FAIL receiver. Quit at  7
	// Simulation took 17.00
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 120*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size
	for ; n > 0; n-- {
		res := p.SendPacket(lib.PACKET_4K, "Server2_1")
		if res == lib.FAIL {
			fmt.Println("FAIL sender. Quit at ", lib.SIM_get_clock())
			p.SIM_wait(10)
			return
		}
	}
	p.SendPacket(lib.PACKET_FINALIZE, "Server2_1")
	fmt.Println("Current time: ", lib.SIM_get_clock(), " sender")
}

func receiver(p *lib.Process, args []string) {
	for {
		packet, res := p.ReceivePacket(args[0])
		if res == lib.FAIL {
			fmt.Println("FAIL receiver. Quit at ", lib.SIM_get_clock())
			return
		}
		if packet == lib.PACKET_FINALIZE {
			break
		}
	}
}

func anomaly(p *lib.Process, args []string) {
	p.CreateHostAnomaly(7, "Server2")
}
