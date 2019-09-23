package main

import (
	"lib"
	"os"
)

func main() {
	lib.SIM_init()
	lib.SIM_platform_init(os.Args[1])

	lib.SIM_function_register("execute1", execute1)
	lib.SIM_function_register("execute2", execute2)
	//lib.SIM_function_register("execute3", execute3)

	lib.SIM_launch_application(os.Args[2])

	lib.SIM_run(nil)
}
