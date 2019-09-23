package main

import (
	"lib"
	"os"
)

func main() {
	lib.SIM_init()
	lib.SIM_platform_init(os.Args[1])

	lib.SIM_function_register("sender", sender)
	lib.SIM_function_register("receiver", receiver)

	lib.SIM_launch_application(os.Args[2])

	lib.SIM_run(nil)
}
