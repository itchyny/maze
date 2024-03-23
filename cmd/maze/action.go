package main

import (
	"fmt"
	"math/rand/v2"
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
			maze.PrintPNG(config.Output, config.Scale)
		case "svg":
			maze.PrintSVG(config.Output, config.Scale)
		}
	}
	return nil
}

func createMaze(config *Config) *maze.Maze {
	maze := maze.NewMaze(config.Height, config.Width,
		maze.WithStart(config.Start),
		maze.WithGoal(config.Goal),
		maze.WithRandSource(rand.NewPCG(config.Seed, config.Seed)),
	)
	maze.Generate()
	if config.Solution {
		maze.Solve()
	}
	return maze
}
