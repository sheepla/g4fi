package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sheepla/g4fi/api"
)

func main() {
	c := api.NewG4fClient(&http.Client{
		Timeout: 30 * time.Second,
	}).WithHost("localhost:8080")

	chat := api.NewChatContext("You", "")
	chat.AddUserMessage(strings.Join(os.Args[1:], " "))

	if err := c.SendAndStreamConversation(chat, os.Stdout); err != nil {
		panic(err)
	}
}
