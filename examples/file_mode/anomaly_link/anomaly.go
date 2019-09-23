package main

import (
	"fmt"
	"lib"
)

func sender(p *lib.Process, args []string) {
	host := p.GetHost()

	task1 := lib.NewTask("task", 100, 5, nil)
	res := p.SendTask(task1, args[0])
	if res == lib.FAIL {
		fmt.Println("Send FAIL")
		lib.SIM_subprocess_create("", waiter, host, nil)
	}
}

func receiver(p *lib.Process, args []string) {
	fmt.Println(lib.SIM_get_clock())
	p.ReceiveTask(args[0])
	fmt.Println(lib.SIM_get_clock())
}

func anomaly(p *lib.Process, args []string) {
	p.CreateLinkAnomaly(3, "link", 0)
}

func waiter(p *lib.Process, args []string) {
	fmt.Println("Starting to wait")
	p.SIM_wait(10)
}
