package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/itchyny/maze"
)

// Config is the command configuration
type Config struct {
	Width       int
	Height      int
	Start       *maze.Point
	Goal        *maze.Point
	Interactive bool
	Image       bool
	Scale       int
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

	width := ctx.GlobalInt("width")
	height := ctx.GlobalInt("height")
	if width <= 0 || height <= 0 {
		termWidth, termHeight, err := terminal.GetSize(0)
		if err != nil {
			return nil, []error{err}
		}
		if width <= 0 {
			width = (termWidth - 4) / 4
		}
		if height <= 0 {
			height = (termHeight - 5) / 2
		}
	}

	start := &maze.Point{X: 0, Y: 0}
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

	goal := &maze.Point{X: height - 1, Y: width - 1}
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

	output := ctx.App.Writer
	outfile := ctx.GlobalString("output")
	if outfile != "" {
		file, err := os.Create(outfile)
		if err != nil {
			errs = append(errs, errors.New("cannot create the output file: "+outfile))
		} else {
			output = file
		}
	}

	image := ctx.GlobalBool("image")
	if image {
		if file, ok := output.(*os.File); ok && isatty.IsTerminal(file.Fd()) {
			errs = append(errs, errors.New("cannot write binary data into the terminal\nuse -output flag"))
		}
	}

	scale := ctx.GlobalInt("scale")

	seed := int64(ctx.GlobalInt("seed"))
	if !ctx.IsSet("seed") {
		seed = time.Now().UnixNano()
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
		Image:       image,
		Scale:       scale,
		Solution:    solution,
		Format:      format,
		Seed:        seed,
		Output:      output,
	}, nil
}
