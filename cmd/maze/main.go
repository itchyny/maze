package main

import (
	"os"
)

var (
	name        = "maze"
	version     = "0.0.8"
	description = "Maze generating and solving program"
	author      = "itchyny"
)

func main() {
	os.Exit(run(os.Args))
}
