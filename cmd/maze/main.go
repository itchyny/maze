package main

import (
	"os"
)

var name = "maze"
var version = "0.0.3"
var description = "Maze generating and solving program"
var author = "itchyny"

func main() {
	os.Exit(run(os.Args))
}
