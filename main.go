package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/sheepla/g4fi/api"
	"github.com/sheepla/g4fi/ui"
)

type Options struct {
	TimeoutSec int    ` arg:"-t, --timeout" help:"Timeout seconds" default:"30"`
	Server     string ` arg:"-s, --server" help:"hostname and port of g4f API instance" default:"localhost:8080"`
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %s", err)
	}

	client := api.NewG4fClient(&http.Client{
		Timeout: time.Duration(opts.TimeoutSec) * time.Second,
	}).WithHost(opts.Server)

	if err := ui.RunInteractiveMode(client); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

func parseArgs(args []string) (*Options, error) {
	var opts Options

	p, err := arg.NewParser(arg.Config{
		Program:   "g4fi",
		IgnoreEnv: false,
	}, &opts)
	if err != nil {
		return &opts, err
	}

	if err := p.Parse(args); err != nil {
		switch {
		case errors.Is(err, arg.ErrHelp):
			p.WriteHelp(os.Stderr)
			os.Exit(1)
		default:
			return &opts, err
		}
	}

	return &opts, nil
}
