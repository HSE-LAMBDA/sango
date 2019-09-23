package main

import (
	"fmt"
	"lib"
)

// 5.5 8.5 16
func sender(p *lib.Process, args []string) {
	task1 := lib.NewTask("task", 100, 5, nil)
	task2 := lib.NewTask("task", 100, 1, nil)
	task3 := lib.NewTask("task", 100, 10, nil)
	p.DetachedSendTask(task1, "B_1")
	p.SIM_wait(2.5)
	fmt.Println("end wait")
	p.DetachedSendTask(task2, "B_2")
	p.DetachedSendTask(task3, "B_3")
}

func receiver(p *lib.Process, args []string) {
	_, _ = p.ReceiveTask(args[0])
	fmt.Println("end listen", lib.SIM_get_clock(), args[0])
}
