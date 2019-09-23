package main

import (
	"fmt"
	"lib"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// 5.5 11 16
func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}

func sender(p *lib.Process, args []string) {
	file := lib.NewFile("1", 10*lib.MB, lib.WRITE, lib.PACKET_4K)
	n := file.Size / lib.PACKET_4K.Size

	for n > 0 {
		p.SendPacket(lib.PACKET_4K, "B_1")
		n--
	}

	p.SendPacket(lib.PACKET_FINALIZE, "B_1")
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
