package src

import (
	lib "gosan"
	"log"
)

type SANVolume struct {
	*NamingProps
	*StorageProperties `json:"props"`
	*SANComponent      `json:"-"`
	disks              map[string]*SANDisk `json:"-"`
}

func initVolumes(iob *IOBalancer) map[string]*SANVolume {
	libVols := lib.GetAllVolumes()
	SANVolumes := make(map[string]*SANVolume)

	for name, vol := range libVols {
		tVol := NewSANVolume(vol)
		SANVolumes[name] = tVol

		//iob.allComponents[name] = tVol
	}
	iob.volumes = SANVolumes
	return SANVolumes
}

func NewUnifiedLogicalVolume(jbods []*SANJBODController) *SANVolume {
	naming := NewNamingProps("AllDisk", "disk-drive", "1")

	capacity := 0.
	avg_write_speed := 0.
	avg_read_speed := 0.

	for _, jbod := range jbods {
		for _, disk := range jbod.disksSlice {
			capacity += disk.RawCapacity
			avg_write_speed += disk.AvgWriteSpeed
			avg_read_speed += disk.AvgReadSpeed
		}
	}

	tv := &SANVolume{
		NamingProps: naming,
		StorageProperties: &StorageProperties{
			FreeSpace:     capacity,
			RawCapacity:   capacity,
			AvgWriteSpeed: avg_write_speed,
			AvgReadSpeed:  avg_read_speed,
			CommonProps: &CommonProps{
				Status: OK,
			},
		},
	}
	return tv
}

func NewSANVolume(libVolume *lib.Volume) *SANVolume {
	naming := NewNamingProps(libVolume.Id, "volume", libVolume.Id)

	capacity := 0.
	avg_write_speed := 0.
	avg_read_speed := 0.

	jbods := lib.GetAllJBODs()

	for _, volumPart := range libVolume.Mounts {
		jbod, ok := jbods[volumPart.Id]
		if !ok {
			log.Panicf("No such jbod %s", volumPart.Id)
		}

		for _, disk := range jbod.DiskArr {
			capacity += disk.Size
			avg_write_speed += disk.WriteRate
			avg_read_speed += disk.ReadRate
		}
	}

	tv := &SANVolume{
		NamingProps: naming,
		SANComponent: &SANComponent{
			currentState: "default",
		},
		StorageProperties: &StorageProperties{
			FreeSpace:     capacity,
			RawCapacity:   capacity,
			AvgWriteSpeed: avg_write_speed,
			AvgReadSpeed:  avg_read_speed,
			CommonProps: &CommonProps{
				Status: OK,
			},
		},
	}

	return tv
}
