package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HttpClient represent standard library http compatible clients.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request stores Discord webhook request information
type Request struct {
	messages []Message // Slice of Discord messages
	url      string    // Discord webhook url
	client   HttpClient
}

// NewRequest create a valid request.
func NewRequest(clt HttpClient, url string, messages []Message) (*Request, error) {
	// Use pointer so we can return nil request, prevents caller to send invalid request.
	req := &Request{messages: divideMessages(messages), url: url, client: clt}
	if err := validateRequest(*req); err != nil {
		return nil, err
	}

	return req, nil
}

// Send request to Discord webhook url via http post. Adjusted to the dynamic rate limit.
func (r *Request) Send() ([]*http.Response, error) {
	var responses []*http.Response
	for _, msg := range r.messages {
		resp, err := makeRequest(msg, r.url, r.client)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if err := respError(resp); err != nil {
			return nil, err
		}

		if err := handleRateLimit(resp.Header); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func makeRequest(msg Message, url string, clt HttpClient) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, formatBody(msg))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	return clt.Do(req)
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
