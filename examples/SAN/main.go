package main

import (
	"gosan/examples/SAN/src"
	"io/ioutil"
	"log"
)

func main() {
	log.SetOutput(ioutil.Discard)
	sf := src.NewSystemFlags()
	sf = src.InitFlags(sf)
	sf = src.ParseFlags(sf)
	shd := src.PlatformInit(sf)
	src.StartSimulation(sf)
	_ = shd
}
