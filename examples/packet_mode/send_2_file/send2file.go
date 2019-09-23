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

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 1*lib.MB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size
	for n > 0 {
		p.SendPacket(lib.PACKET_4K, args[0])
		n--
	}
	p.SendPacket(lib.PACKET_FINALIZE, args[0])
	fmt.Println("Current time: ", lib.SIM_get_clock(), " sender")
}

func receiver(p *lib.Process, args []string) {
	for {
		packet, _ := p.ReceivePacket(args[0])
		if packet == lib.PACKET_FINALIZE {
			return
		}
	}
	fmt.Println("Current time: ", lib.SIM_get_clock(), " receiver")
}
