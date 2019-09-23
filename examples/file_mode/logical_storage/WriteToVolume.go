package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

var tatlin Tatlin

type Tatlin struct {
	volume *lib.Volume
}

func StorageWrite(p *lib.Process, args []string) {
	InitVolumes()

	task1 := lib.NewTask("task", 100, 500, nil)
	p.WriteToVolume(tatlin.volume, task1)
	fmt.Println(lib.SIM_get_clock())
}

func InitVolumes() {
	storages := lib.GetDiskDrives()
	volumeNames := make([]string, len(storages))

	i := 0
	for k := range storages {
		volumeNames[i] = k
		i++
	}

	volume := lib.NewVolume(volumeNames[0], storages[volumeNames[0]])
	tatlin.volume = volume
}
