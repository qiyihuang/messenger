package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Author represents author object in an embed object.
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Field represents field object in an embed object.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer represents footer object in an embed object.
type Footer struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Embed represents an embed object in message object.
type Embed struct {
	Author      Author  `json:"author,omitempty"`
	Color       int     `json:"color,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	URL         string  `json:"url,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Footer      Footer  `json:"footer,omitempty"`
}

// Message represents a webhook message.
type Message struct {
	Username string  `json:"username,omitempty"`
	Embeds   []Embed `json:"embeds,omitempty"`
	Content  string  `json:"content,omitempty"`
}

func formatBody(msg Message) (io.Reader, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(jsonMsg)
	return body, nil
}

// MakeRequest make http request to Discord server to send the message.
func MakeRequest(msg Message, url string) (*http.Response, error) {
	body, err := formatBody(msg)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	// Log only when Discord API message (usually error) exists.
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if message, ok := respBody["message"]; ok {
		errMsg := "Discord API error: " + fmt.Sprintf("%v", message)
		return nil, errors.New(errMsg)
	}

	return resp, nil
}
