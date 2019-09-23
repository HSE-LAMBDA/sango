package lib

import (
	"fmt"
)

/*
===========================

Logging

===========================
*/

func (process *Process) GetHost() *Host {
	return process.host
}

func (host *Host) GetName() string {
	return host.Name
}

func (host *Host) GetSpeed() float64 {
	return host.Speed
}

func (host *Host) GetAmountOfTasks() uint64 {
	number := 0
	for _, core := range host.coreManager.coreMap {
		number += len(core.eQueue)
	}
	return uint64(number)
}

func (host *Host) GetType() string {
	return host.Type
}

func (host *Host) GetDevTemp() float64 {
	return 22
}

func getHostByNameFactory(hostsMap map[string]*Host) func(string) *Host {
	return func(hostName string) *Host {
		host, ok := hostsMap[hostName]
		if !ok {
			panic(fmt.Sprintf("No such host =( %s", hostName))
		}
		return host
	}
}

func getHostsFactory(hosts map[string]*Host) func() map[string]*Host {
	return func() map[string]*Host {
		return hosts
	}
}
