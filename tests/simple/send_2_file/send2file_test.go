package main

import (
	"fmt"
	"lib"
	"strconv"
)

// 5.5 11 16

func ExampleSendTwoFilesBetweenTwoHosts() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)

	lib.SIM_launch_application("topology/deployment.xml")

	lib.SIM_run(nil)

	// Output:
	// End of transferring at B_1. Time is 20.00
	// End of transferring at B_2. Time is 35.00
	// Simulation took 35.00
}

func sender(p *lib.Process, args []string) {
	size, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Arguments are not correct")
	}
	file := lib.NewFile("1", size*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size
	for n > 0 {
		p.SendPacket(lib.PACKET_4K, args[1])
		n--
	}
	p.SendPacket(lib.PACKET_FINALIZE, args[1])
	fmt.Printf("End of transferring at %s. Time is %.2f\n", args[1], lib.SIM_get_clock())
}

func receiver(p *lib.Process, args []string) {
	for {
		packet, _ := p.ReceivePacket(args[0])
		if packet == lib.PACKET_FINALIZE {
			return
		}
	}
}
