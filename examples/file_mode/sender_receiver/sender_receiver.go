package main

import (
	"fmt"
	"lib"
)

// 5.5 11 16
func sender(p *lib.Process, args []string) {
	task1 := lib.NewTask("task", 100, 5, nil)
	p.SendTask(task1, "B_1")
	p.SIM_wait(10)
	fmt.Println("Current time: ", lib.SIM_get_clock(), " sender")
}

func receiver(p *lib.Process, args []string) {
	p.ReceiveTask(args[0])
	fmt.Println("Receive at time: ", lib.SIM_get_clock())
	p.SIM_wait(15)
	fmt.Println("Current time: ", lib.SIM_get_clock(), " receiver")

}
