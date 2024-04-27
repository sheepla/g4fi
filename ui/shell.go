package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"github.com/sheepla/g4fi/api"
)

const helpMessage = `
Usage:
	/help         Show help message
	/quit, Ctrl-D Quit interactive session
`

type State struct {
	AiProvider string
	AiModel    string
	Messages   []api.Message
}

func initState() *State {

	return &State{
		AiProvider: "",
		AiModel:    "",
		Messages:   []api.Message{},
	}
}

func (s *State) AddUserMessage(message string) {
	s.Messages = append(s.Messages, api.Message{
		Role:    "user",
		Content: message,
	})
}

func (s *State) AddAiMessage(message string) {
	s.Messages = append(s.Messages, api.Message{
		Role:    "assistant",
		Content: message,
	})
}

//func (s *State) ToPromptText() string {
//	provider := s.AiProvider
//	if provider == "" {
//		provider = "Auto"
//	}
//
//	model := s.AiModel
//	if model == "" {
//		model = "Default"
//	}
//
//	return fmt.Sprintf("\n[ %s %s ]", provider, model)
//}

func initCompleter() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("/help"),
		readline.PcItem("/quit"),
	)
}

func RunInteractiveMode(client *api.G4fClient) error {
	state := initState()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "",
		AutoComplete:    initCompleter(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to initialize readline module")
	}
	defer rl.Close()

LOOP:
	for {
		//fmt.Println(state.ToPromptText())

		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		switch strings.TrimSpace(line) {
		case "/help":
			fmt.Println(helpMessage)
			continue LOOP
		case "/quit":
			return nil
		case "":
			continue LOOP
		default:
			if err := invokeAi(client, state); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			state.AddUserMessage(line)
		}

	}
}

func invokeAi(client *api.G4fClient, state *State) error {
	return client.SendAndStreamConversation(&api.Conversation{
		Provider: state.AiProvider,
		Model:    state.AiModel,
		Messages: state.Messages,
	}, os.Stdout)
}
