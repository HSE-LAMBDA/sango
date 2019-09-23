package main

import (
	"fmt"
	"lib"
	"os"
)

// 5.5 11 16

func main() {
	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("execute", execute)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func execute(p *lib.Process, args []string) {
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_4K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_8K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_16K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_32K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_64K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_128K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_256K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_512K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_1024K)
	fmt.Println(lib.SIM_get_clock())
	p.Execute(lib.PACKET_2048K)
	fmt.Println(lib.SIM_get_clock())
}
