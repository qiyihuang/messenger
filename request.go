package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/qiyihuang/messenger/ratelimit"
)

var post = http.Post

// Request stores Discord webhook request information
type Request struct {
	Messages []Message // Slice of Discord messages
	URL      string    // Discord webhook url
}

// Send sends the request to Discord webhook url via http post. Request is
// validated and send speed adjusted by rate limiter.
func (r Request) Send() ([]*http.Response, error) {
	r.Messages = divideMessages(r.Messages)
	if err := validateRequest(r); err != nil {
		return nil, err
	}

	var responses []*http.Response
	for _, msg := range r.Messages {
		resp, err := post(r.URL, "application/json", formatBody(msg))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if err := respError(resp); err != nil {
			return nil, err
		}

		if err := ratelimit.Wait(resp.Header); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func countEmbed(e Embed) int {
	total := len(e.Title) + len(e.Description) + len(e.Author.Name) + len(e.Footer.Text)
	for _, field := range e.Fields {
		total += len(field.Name)
		total += len(field.Value)
	}
	return total
}

func divideEmbeds(msg Message) (dividedEmbeds [][]Embed) {
	var total int
	var startIndex int
	for i, e := range msg.Embeds {
		count := countEmbed(e)
		total += count
		// i will be included in next chunk.
		if total > EmbedTotalLimit || i-startIndex == MessageEmbedNumLimit {
			dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:i])
			startIndex = i
			// Set current count to initial total of next message.
			total = count
		}
	}
	// Add the last chunk.
	dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:])
	return
}

// divideMessages breaks message into multiple messages depending on embed total
// character count and number of embeds.
func divideMessages(messages []Message) (msgs []Message) {
	for _, msg := range messages {
		dividedEmbeds := divideEmbeds(msg)
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
	return bytes.NewBuffer(jsonMsg)
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

	// Discord API error message is written in "message" field in response body.
	// "message" is of string type. "https://discord.com/developers/docs/reference#error-messages"
	if message, ok := respBody["message"]; ok {
		return errors.New("Discord API error: " + message.(string))
	}

	return nil
}
