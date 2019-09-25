package deepcontroller

import (
	lib "gosan"
)


func Init_d(flags *lib.DeepControllerFlags) *Manager  {

	// Deep controller
	// lib.SIM_subprocess_create_async(name, PIDecorator(f), host, NewSANProcess(), data)

	host := lib.GetHostByName("Helper")


	// Deep controller
	// lib.SIM_subprocess_create_async(name, PIDecorator(f), host, NewSANProcess(), data)
	dcm := NewManager(flags)
	p := lib.NewProcess("DC", host)
	lib.ProcWrapperTemp(dcm.DeepControllerManagerProcess, p, flags)

	return dcm
}
