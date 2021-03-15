package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Msg Message
	URL string
}

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

// send sends the message to Discord via http.
func (r Request) send() (resp *http.Response, err error) {
	err = validateURL(r.URL)
	if err != nil {
		return
	}

	err = validateMessage(r.Msg)
	if err != nil {
		return
	}

	body := formatBody(r.Msg)
	resp, err = http.Post(r.URL, "application/json", body)
	if err != nil {
		return
	}

	err = respError(resp)
	return
}
