package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func formatBody(msg Message) io.Reader {
	// Marshal would never fail since Discord webhook message does not
	// contain types not supported by Marshal.
	jsonMsg, _ := json.Marshal(msg)
	body := bytes.NewBuffer(jsonMsg)
	return body
}

func respError(resp *http.Response) error {
	// Log only when Discord API message (usually error) exists.
	var respBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	if message, ok := respBody["message"]; ok {
		errMsg := "Discord API error: " + fmt.Sprintf("%v", message)
		return errors.New(errMsg)
	}

	return nil
}

// makeRequest sends the message to Discord via http.
func makeRequest(msg Message, url string) (resp *http.Response, err error) {
	body := formatBody(msg)
	resp, err = http.Post(url, "application/json", body)
	if err != nil {
		return
	}

	err = respError(resp)
	return
}
