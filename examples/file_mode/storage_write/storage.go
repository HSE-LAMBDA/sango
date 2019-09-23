package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

// 5.5 11 16

func StorageWrite(p *lib.Process, args []string) {
	disk := p.GetHost().GetStorage()
	task1 := lib.NewTask("task", 100, 5, nil)
	task2 := lib.NewTask("task", 100, 1, nil)
	task3 := lib.NewTask("task", 100, 10, nil)
	p.WriteAsync(disk, task1)
	fmt.Println(lib.SIM_get_clock())
	p.SIM_wait(2.5)
	p.WriteAsync(disk, task3)
	p.SIM_wait(2.5)
	p.WriteAsync(disk, task2)
}
