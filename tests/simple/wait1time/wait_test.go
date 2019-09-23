package main

import (
	"fmt"
	"lib"
)

func ExampleOfWaiting1Time() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")
	lib.SIM_function_register("wait", wait)
	lib.SIM_launch_application("topology/deployment.xml")

	lib.SIM_run(nil)

	// Output:
	// Current time: 10.00
	// Simulation took 10.00
}

func wait(p *lib.Process, args []string) {
	res := p.SIM_wait(10)
	if res != lib.OK {
		panic("ERROR")
	}
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
}
