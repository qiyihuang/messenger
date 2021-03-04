package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/qiyihuang/messenger/pkg/message"
)

func formatBody(msg message.Message) (io.Reader, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(jsonMsg)
	return body, nil
}

func respError(resp *http.Response) error {
	// Log only when Discord API message (usually error) exists.
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if message, ok := respBody["message"]; ok {
		errMsg := "Discord API error: " + fmt.Sprintf("%v", message)
		return errors.New(errMsg)
	}

	return nil
}

// Send sends the message to Discord via http.
func Send(msg message.Message, url string) (resp *http.Response, err error) {
	body, err := formatBody(msg)
	if err != nil {
		return
	}

	resp, err = http.Post(url, "application/json", body)
	if err != nil {
		return
	}

	err = respError(resp)
	return
}
