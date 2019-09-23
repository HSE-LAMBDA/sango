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

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)
	lib.SIM_function_register("anomaly", anomaly)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 120*lib.KB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size
	for ; n > 0; n-- {
		res := p.SendPacket(lib.PACKET_4K, "B_1")
		if res == lib.FAIL {
			fmt.Println("FAIL sender. Quit at ", lib.SIM_get_clock())
			return
		}
	}
	p.SendPacket(lib.PACKET_FINALIZE, "B_1")
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
	fmt.Println("Current time: ", lib.SIM_get_clock(), " receiver")
}

func anomaly(p *lib.Process, args []string) {
	p.CreateHostAnomaly(10, "B")
}
