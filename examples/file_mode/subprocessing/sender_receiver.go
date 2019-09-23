package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

func launcher(p *lib.Process, args []string) {
	host := p.GetHost()
	for i := 0; i < 2; i++ {
		lib.SIM_subprocess_create("", sender, host, nil)
	}
}

func sender(p *lib.Process, args []string) {
	task1 := lib.NewTask("task", 100, 5, nil)
	p.SendTask(task1, "1")
}

func receiver(p *lib.Process, args []string) {
	for i := 0; i < 2; i++ {
		fmt.Println("start listen", lib.SIM_get_clock())
		_ = p.ReceiveTask(args[0])
		fmt.Println("end listen", lib.SIM_get_clock())
	}
}
