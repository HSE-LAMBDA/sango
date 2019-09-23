package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Println

func sender(p *lib.Process, args []string) {
	// task1 := lib.NewTask("task", 100, 5, nil)
	// fmt.Println(lib.SIM_get_clock())
	// p.SendTask(task1, "1")
	// fmt.Println(lib.SIM_get_clock())
}

func receiver(p *lib.Process, args []string) {
	fmt.Println("start listen", lib.SIM_get_clock())
	_, res := p.ReceiveTaskWithTimeout(args[0], 3.45)
	fmt.Println("end listen", lib.SIM_get_clock(), res)
	p.SIM_wait(10)
	_, res = p.ReceiveTaskWithTimeout(args[0], 3.45)
	fmt.Println("end listen", lib.SIM_get_clock(), res)
}
