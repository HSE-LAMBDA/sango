package main

import (
	"fmt"
	"lib"
)

func sender(p *lib.Process, args []string) {
	task1 := lib.NewTask("task", 100, 5, nil)
	p.SendTask(task1, "1")
}

func receiver(p *lib.Process, args []string) {
	fmt.Println("start listen", lib.SIM_get_clock())
	_, _ = p.ReceiveTask(args[0])
	fmt.Println("end listen", lib.SIM_get_clock())

}
