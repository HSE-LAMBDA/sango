package main

import (
	"fmt"
	"lib"
	"os"
)

func main() {
	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 100*lib.MB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size

	for ; n > 0; n-- {
		p.SendPacket(lib.PACKET_4K, "B_1")
		p.ReceivePacket("A_1")
	}
	p.SendPacket(lib.PACKET_FINALIZE, "B_1")
	fmt.Println("Current time: ", lib.SIM_get_clock(), " sender")
}

func receiver(p *lib.Process, args []string) {
	defer fmt.Println("Current time: ", lib.SIM_get_clock(), " receiver")
	storage := p.GetHost().GetOneStorage()

	for {
		packet, _ := p.ReceivePacket("B_1")
		if packet == lib.PACKET_FINALIZE {
			return
		}

		res := p.WriteSync(storage, packet)
		if res == lib.FAIL {
			return
		}
		p.SendPacket(lib.PACKET_ACK, "A_1")

	}

}
