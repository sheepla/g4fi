package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	urlbuilder "github.com/sheepla/go-urlbuilder"
	"github.com/tidwall/gjson"
)

type G4fClient struct {
	internal *http.Client
	useSSL   bool
	host     string
}

func (c *G4fClient) BaseURL() string {
	if c.useSSL {
		return fmt.Sprintf("https://%s", c.host)
	} else {
		return fmt.Sprintf("http://%s", c.host)
	}
}

func NewG4fClient(c *http.Client) *G4fClient {
	return &G4fClient{
		internal: c,
		useSSL:   false,
		host:     "localhost:8080",
	}
}

func (c *G4fClient) WithHost(host string) *G4fClient {
	c.host = host
	return c
}

func (c *G4fClient) UseSSL(u bool) *G4fClient {
	c.useSSL = u
	return c
}

type Providers []string

func (c *G4fClient) GetProviders() (*Providers, error) {
	url, err := urlbuilder.MustParse(c.BaseURL()).
		SetPath("/backend-api/v2/providers").String()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.internal.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var providers Providers
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, err
	}

	return &providers, nil
}

type Models []string

func (c *G4fClient) GetModels(provider string) (*Models, error) {
	url, err := urlbuilder.MustParse(c.BaseURL()).
		SetPath("/backend-api/v2/models").AppendPath(provider).String()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.internal.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var models Models
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	return &models, nil
}

type ChatContext struct {
	//ID             string `json:"id"`
	//ConversationID string `json:"conversation_id"`
	//Jailbreak      string `json:"jailbreak"`
	//WebSearch      bool   `json:"web_search"`
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Provider string    `json:"provider"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatContext(provider string, modelName string) *ChatContext {
	return &ChatContext{
		Messages: []Message{},
		Provider: provider,
		Model:    modelName,
	}
}

func (chat *ChatContext) AddUserMessage(message string) {
	chat.Messages = append(chat.Messages, Message{
		Role:    "user",
		Content: message,
	})
}

func (chat *ChatContext) AddAssistantMessage(message string) {
	chat.Messages = append(chat.Messages, Message{
		Role:    "assistant",
		Content: message,
	})
}

func (c *G4fClient) sendChatConversation(
	chat *ChatContext,
	handler func(r io.Reader) error,
) error {
	url, err := urlbuilder.MustParse(c.BaseURL()).
		SetPath("/backend-api/v2/conversation").String()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.internal.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()

	if err := handler(resp.Body); err != nil {
		return err
	}

	return nil
}

// Reply responses return as a stream where each line becomes JSON
// so this function extracts message tokens line by line from the io.Reader
// and writes them to an destination io.Writer.
func extractReplyMessageStream(src io.Reader, dest io.Writer) error {
	sc := bufio.NewScanner(src)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			return err
		}

		res := gjson.ParseBytes(sc.Bytes())

		if res.Get("type").String() == "content" {
			fmt.Fprint(dest, res.Get("content").String())
		}
	}

	return nil
}

func (c *G4fClient) SendAndStreamConversation(
	chat *ChatContext,
	dest io.Writer,
) error {
	err := c.sendChatConversation(chat, func(r io.Reader) error {
		return extractReplyMessageStream(r, dest)
	})
	if err != nil {
		return err
	}

	return nil
}
