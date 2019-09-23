package main

import (
	"lib"
)

func ExampleExecute10() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("execute", execute)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Simulation took 10.00
}

func execute(p *lib.Process, args []string) {
	res := p.Execute(lib.PACKET_4K)
	if res != lib.OK {
		panic("ERROR")
	}
}
