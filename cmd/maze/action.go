package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"

	"github.com/itchyny/maze"
)

func action(ctx *cli.Context) error {
	config, errors := makeConfig(ctx)
	if errors != nil {
		hasErr := false
		for _, err := range errors {
			if err.Error() != "" {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
				hasErr = true
			}
		}
		if hasErr {
			fmt.Fprint(os.Stderr, "\n")
		}
		cli.ShowAppHelp(ctx)
		return nil
	}

	maze := createMaze(config)
	if config.Interactive {
		err := termbox.Init()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return nil
		}
		defer termbox.Close()
		interactive(maze, config.Format)
	} else {
		switch config.Image {
		case "":
			maze.Print(config.Output, config.Format)
		case "png":
			maze.PrintPNG(config.Output, config.Format, config.Scale)
		case "svg":
			maze.PrintSVG(config.Output, config.Format, config.Scale)
		default:
			return fmt.Errorf("unsupported image format: %s", config.Image)
		}
	}
	return nil
}

func createMaze(config *Config) *maze.Maze {
	//lint:ignore SA1019 Random seed is necessary for testing.
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
