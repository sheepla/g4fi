package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sheepla/g4fi/api"
	"github.com/sheepla/g4fi/ui"
)

func main() {
	client := api.NewG4fClient(&http.Client{
		Timeout: 30 * time.Second,
	})
	if err := ui.RunInteractiveMode(client); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s", err)
	}
}
