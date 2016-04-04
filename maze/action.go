package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/codegangsta/cli"
	"github.com/itchyny/maze"
	"github.com/nsf/termbox-go"
)

func action(ctx *cli.Context) {
	err := termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	config, errors := makeConfig(ctx)
	if errors != nil {
		termbox.Close()
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		cli.ShowAppHelp(ctx)
		return
	}

	maze := createMaze(config)
	if config.Interactive {
		defer termbox.Close()
		interactive(maze, config.Format)
	} else {
		termbox.Close()
		maze.Print(config.Output, config.Format)
	}
}

func createMaze(config *Config) *maze.Maze {
	rand.Seed(config.Seed)
	maze := maze.NewMaze(config.Height, config.Width)
	maze.Start = config.Start
	maze.Goal = config.Goal
	maze.Cursor = config.Start
	maze.Generate()
	if config.Solution {
		maze.Solve()
	}
	return maze
}
