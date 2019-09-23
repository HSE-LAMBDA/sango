package main

import (
	"lib"
)

func ExampleExecute5Packet() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("execute", execute)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Simulation took 50.00
}

func execute(p *lib.Process, args []string) {
	res1 := p.Execute(lib.PACKET_4K)
	if res1 != lib.OK {
		panic("ERROR")
	}
	res2 := p.Execute(lib.PACKET_4K)
	if res2 != lib.OK {
		panic("ERROR")
	}
	res3 := p.Execute(lib.PACKET_4K)
	if res3 != lib.OK {
		panic("ERROR")
	}
	res4 := p.Execute(lib.PACKET_4K)
	if res4 != lib.OK {
		panic("ERROR")
	}
	res5 := p.Execute(lib.PACKET_4K)
	if res5 != lib.OK {
		panic("ERROR")
	}
}
