package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Printf

func sender(p *lib.Process, args []string) {
	task := lib.NewTask("task", 100, 5, nil)
	x := p.Execute(task)
	fmt.Println(x)
	fmt.Println("end of wait ", lib.SIM_get_clock())
	task1 := lib.NewTask("task", 100, 5, nil)

	res := p.SendTask(task1, args[0])
	if res == lib.FAIL {
		fmt.Println("Send FAIL")
	}

}

func receiver(p *lib.Process, args []string) {
	fmt.Println("Start:  ", lib.SIM_get_clock())
	_, res := p.ReceiveTask(args[0])
	if res == lib.FAIL {
		fmt.Println("receive FAIL")
	}
	fmt.Println("End:  ", lib.SIM_get_clock())
}

func anomaly(p *lib.Process, args []string) {
	hostId := "A"
	p.CreateHostAnomaly(3, hostId)
}
