package main

import (
	"lib"
	"os"
)

func main() {
	lib.SIM_init(os.Args[1])
	lib.SIM_platform_init(os.Args[2])

	lib.SIM_function_register("storage", StorageWrite)

	lib.SIM_launch_application(os.Args[3])

	lib.SIM_run(nil)
}
