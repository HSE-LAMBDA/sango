package main

import (
	"fmt"
	"lib"
)

func ExampleSendOneFileBetweenTwoHostsDetached() {

	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)

	lib.SIM_launch_application("topology/deployment.xml")

	lib.SIM_run(nil)
	// Output:
	// Current time: 1.00
	// Simulation took 1.00
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 4*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size
	for n > 0 {
		p.DetachedSendPacket(lib.PACKET_4K, "B_1")
		n--
	}
	p.SendPacket(lib.PACKET_FINALIZE, "B_1")
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
}

func receiver(p *lib.Process, args []string) {
	for {
		packet, _ := p.ReceivePacket(args[0])
		if packet == lib.PACKET_FINALIZE {
			return
		}
	}
}
