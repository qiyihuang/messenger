package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/qiyihuang/messenger/pkg/ratelimit"
)

// Request stores Discord webhook request information
type Request struct {
	Messages []Message // Slice of Discord messages
	URL      string    // Discord webhook url
}

var post = http.Post

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
	switch {
	// Body is empty.
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	}

	//Discord API error message is written in "message" field in response body.
	if message, ok := respBody["message"]; ok {
		errMsg := "Discord API error: " + fmt.Sprintf("%v", message)
		return errors.New(errMsg)
	}

	return nil
}

// Send sends the request to Discord webhook url via http post. Request is
// validated and send speed adjusted by rate limiter.
func (r Request) Send() (responses []*http.Response, err error) {
	err = validateRequest(r)
	if err != nil {
		return
	}

	for _, msg := range r.Messages {
		body := formatBody(msg)
		var resp *http.Response
		resp, err = post(r.URL, "application/json", body)
		if err != nil {
			return
		}

		defer resp.Body.Close()
		responses = append(responses, resp)

		err = respError(resp)
		if err != nil {
			return
		}

		err = ratelimit.Wait(resp.Header)
		if err != nil {
			return
		}
	}

	return
}
