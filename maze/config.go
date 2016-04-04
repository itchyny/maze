package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/itchyny/maze"
	"github.com/nsf/termbox-go"
)

// Config is the command configuration
type Config struct {
	Width       int
	Height      int
	Start       *maze.Point
	Goal        *maze.Point
	Interactive bool
	Solution    bool
	Format      *maze.Format
	Seed        int64
	Output      io.Writer
}

func makeConfig(ctx *cli.Context) (*Config, []error) {
	var errs []error

	if ctx.GlobalBool("help") {
		errs = append(errs, errors.New(""))
		return nil, errs
	}

	termWidth, termHeight := termbox.Size()

	width := ctx.GlobalInt("width")
	if width <= 0 {
		width = (termWidth - 4) / 4
	}

	height := ctx.GlobalInt("height")
	if height <= 0 {
		height = (termHeight - 5) / 2
	}

	start := &maze.Point{0, 0}
	starts := strings.Split(ctx.GlobalString("start"), ",")
	if len(starts) > 0 {
		if value, err := strconv.Atoi(starts[0]); err == nil {
			if 0 <= value && value < height {
				start.X = value
			}
		}
	}
	if len(starts) > 1 {
		if value, err := strconv.Atoi(starts[1]); err == nil {
			if 0 <= value && value < width {
				start.Y = value
			}
		}
	}

	goal := &maze.Point{height - 1, width - 1}
	goals := strings.Split(ctx.GlobalString("goal"), ",")
	if len(goals) > 0 {
		if value, err := strconv.Atoi(goals[0]); err == nil {
			if 0 <= value && value < height {
				goal.X = value
			}
		}
	}
	if len(goals) > 1 {
		if value, err := strconv.Atoi(goals[1]); err == nil {
			if 0 <= value && value < width {
				goal.Y = value
			}
		}
	}

	interactive := ctx.GlobalBool("interactive")

	solution := ctx.GlobalBool("solution")

	format := maze.Default
	if ctx.GlobalString("format") == "color" {
		format = maze.Color
	}

	seed := int64(ctx.GlobalInt("seed"))
	if !ctx.IsSet("seed") {
		seed = time.Now().UnixNano()
	}

	output := ctx.App.Writer
	outfile := ctx.GlobalString("output")
	if outfile != "" {
		file, err := os.Create(outfile)
		if err != nil {
			errs = append(errs, errors.New("Cannot create the output file: "+outfile+"\n\n"))
		} else {
			output = file
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Config{
		Width:       width,
		Height:      height,
		Start:       start,
		Goal:        goal,
		Interactive: interactive,
		Solution:    solution,
		Format:      format,
		Seed:        seed,
		Output:      output,
	}, nil
}
