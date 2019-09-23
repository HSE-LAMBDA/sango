package main

import (
	"fmt"
	"lib"
)

var _ = fmt.Print

func anomaly(p *lib.Process, args []string) {
	p.CreateLinkAnomaly(0, "link", 0)
}
