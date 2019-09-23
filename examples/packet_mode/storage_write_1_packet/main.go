package main

import (
	"fmt"
	"lib"
	"os"
)

func main() {
	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("write_packet", write_1_packet)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func write_1_packet(p *lib.Process, args []string) {
	storage := p.GetHost().GetOneStorage()
	p.WriteSync(storage, lib.PACKET_4K)
	fmt.Println("Current time: ", lib.SIM_get_clock(), " receiver")
}
