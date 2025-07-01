package mattermost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"financing-offer/pkg/environment"
)

type Client struct {
	client     http.Client
	webhookUrl string
	env        environment.Environment
}

func NewClient(env environment.Environment, webhookUrl string) *Client {
	return &Client{
		client:     http.Client{},
		webhookUrl: webhookUrl,
		env:        env,
	}
}

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (c *Client) Send(channel, message string) error {
	if !c.env.IsProduction() {
		channel += "-uat"
	}
	messageBytes, _ := json.Marshal(
		Message{
			Channel: channel,
			Text:    message,
		},
	)
	req, err := http.NewRequest(http.MethodPost, c.webhookUrl, bytes.NewReader(messageBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = c.client.Do(req)
	if err != nil {
		return fmt.Errorf("matter most client Send %w", err)
	}
	return nil
}
