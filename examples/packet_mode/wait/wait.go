package main

import (
	"fmt"
	"lib"
	"os"
)

// 5.5 11 16

func main() {
	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("wait", wait)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(100.)
}

func wait(tatlin *lib.Process, args []string) {
	tatlin.SIM_wait(10)
	fmt.Println(lib.SIM_get_clock())
}
