package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Request stores Discord webhook request information
type Request struct {
	Msg Message
	URL string
}

// httpPoster sends http POST requests. e.g. http.Client
type httpPoster interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// formatBody serialises Message.
func formatBody(msg Message) io.Reader {
	// Marshal would never fail since Discord webhook message does not
	// contain types not supported by Marshal.
	jsonMsg, _ := json.Marshal(msg)
	body := bytes.NewBuffer(jsonMsg)
	return body
}

func respError(resp *http.Response) error {
	var respBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	//Discord API error message is written in "message" field in response body.
	if message, ok := respBody["message"]; ok {
		errMsg := "Discord API error: " + fmt.Sprintf("%v", message)
		return errors.New(errMsg)
	}

	return nil
}

// send sends the message to Discord via http.
func (r Request) send(p httpPoster) (resp *http.Response, err error) {
	err = validateURL(r.URL)
	if err != nil {
		return
	}

	err = validateMessage(r.Msg)
	if err != nil {
		return
	}

	body := formatBody(r.Msg)
	resp, err = p.Post(r.URL, "application/json", body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = respError(resp)
	return
}
