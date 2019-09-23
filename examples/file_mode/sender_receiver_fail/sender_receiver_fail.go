package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

func sender(p *lib.Process, args []string) {
	p.SIM_wait(10)
	task1 := lib.NewTask("task", 100, 5, nil)
	res := p.SendTask(task1, "1")
	fmt.Println("Send", "STATUS:", res, "time:", lib.SIM_get_clock())
}

func receiver(p *lib.Process, args []string) {
	fmt.Println("start listen", lib.SIM_get_clock())
	_, res := p.ReceiveTask(args[0])
	fmt.Println("Receiver", "STATUS", res)
	fmt.Println("end listen", lib.SIM_get_clock())
}
