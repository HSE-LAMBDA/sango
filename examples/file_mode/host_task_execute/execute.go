package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

func execute1(p *lib.Process, args []string) {
	task := lib.NewTask("task", 100, 5, nil)
	p.SIM_wait(50)
	p.SIM_wait(50)
	p.SIM_wait(50)
	p.SIM_wait(50)
	p.Execute(task)
	fmt.Println(p.GetName(), lib.SIM_get_clock())
}

func execute2(p *lib.Process, args []string) {
	task := lib.NewTask("task", 100, 5, nil)
	p.Execute(task)
	fmt.Println(p.GetName(), lib.SIM_get_clock())
}

func execute3(p *lib.Process, args []string) {
	//task := lib.NewTask("task", 100, 5, nil)
	//p.Execute(task)
	fmt.Println(p.GetName(), lib.SIM_get_clock())
}
