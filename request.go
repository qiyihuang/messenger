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

func countEmbed(e Embed) int16 {
	total := len(e.Title) + len(e.Description) + len(e.Author.Name) + len(e.Footer.Text)
	for _, field := range e.Fields {
		total += len(field.Name)
		total += len(field.Value)
	}
	int16Total := int16(total)
	return int16Total
}

// divideMessages breaks message into multiple messages depending on embed total
// character count and number of embeds.
func divideMessages(messages []Message) (msgs []Message) {
	for _, msg := range messages {
		var total int16
		var dividedEmbeds [][]Embed
		startIndex := 0
		for i, e := range msg.Embeds {
			count := countEmbed(e)
			total += count
			if total > EmbedTotalLimit {
				dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:i])
				startIndex = i
				total = count
			}
		}
		// Add the last chunk.
		dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:])

		// Create message for every embed chunk.
		for i, embeds := range dividedEmbeds {
			// First message contains content from original message.
			var content string
			if i == 0 {
				content = msg.Content
			}

			msgs = append(msgs, Message{Username: msg.Username, Embeds: embeds, Content: content})
		}
	}

	return msgs
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

	r.Messages = divideMessages(r.Messages)
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
