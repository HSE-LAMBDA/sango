package main

import (
	"fmt"
	"lib"
)

// 5.5 11 16

func ExampleOfWaitingFourTimes() {
	lib.SIM_init("topology/packet_config.json")
	lib.SIM_platform_init("topology/platform.xml")

	lib.SIM_function_register("wait", wait)

	lib.SIM_launch_application("topology/deployment.xml")
	lib.SIM_run(nil)

	// Output:
	// Current time: 10.00
	// Current time: 30.00
	// Current time: 60.00
	// Current time: 100.00
	// Simulation took 100.00
}

func wait(p *lib.Process, args []string) {
	p.SIM_wait(10)
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
	p.SIM_wait(20)
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
	p.SIM_wait(30)
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
	p.SIM_wait(40)
	fmt.Printf("Current time: %.2f\n", lib.SIM_get_clock())
}
