package main

import (
	lib "gosan"
	"gosan/examples/SAN/src"
	"io/ioutil"
	"log"
)

func main() {
	log.SetOutput(ioutil.Discard)
	sf := lib.NewSystemFlags()
	sf = lib.InitFlags(sf)
	sf = lib.ParseFlags(sf)
	shd := src.PlatformInit(sf)
	src.StartSimulation(sf)
	_ = shd
}
