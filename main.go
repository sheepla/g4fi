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

type exitStatus int

const (
	exitStatusOK exitStatus = iota
	exitStatusErrArgs
	exitStatusErrApi
	exitStatusErrInteractiveSession
)

func (e exitStatus) Int() int {
	return int(e)
}

type Options struct {
	TimeoutSec      int    ` arg:"-t, --timeout" help:"Timeout seconds" default:"30"`
	Server          string ` arg:"-s, --server" help:"hostname and port of g4f API instance" default:"localhost:8080"`
	AiProvider      string `arg:"-p, --provider" help:"AI assistant provider name" default:""`
	AiModel         string `arg:"-m, --model" help:"AI assistant model name" default:""`
	ShowAiProviders bool   `arg:"-P, --show-providers" help:"Show a list of AI assistant providers"`
	ShowAiModels    string `arg:"-M, --show-models" help:"Show a list of AI model of the specified provider"`
}

func main() {
	exitStatus := run(os.Args[1:])
	os.Exit(exitStatus.Int())
}

func run(args []string) exitStatus {
	opts, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %s", err)
		return exitStatusErrArgs
	}

	client := api.NewG4fClient(&http.Client{
		Timeout: time.Duration(opts.TimeoutSec) * time.Second,
	}).WithHost(opts.Server)

	if opts.ShowAiProviders {
		providers, err := client.GetProviders()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitStatusErrApi
		}

		for _, p := range *providers {
			fmt.Println(p)
		}

		return exitStatusOK
	}

	if opts.ShowAiModels != "" {
		models, err := client.GetModels(opts.ShowAiModels)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitStatusErrApi
		}

		for _, m := range *models {
			fmt.Println(m)
		}

		return exitStatusOK
	}

	if err := ui.RunInteractiveMode(client, opts.AiProvider, opts.AiModel); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return exitStatusErrInteractiveSession
	}

	return exitStatusOK
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
