package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// HttpClient represent standard library http compatible clients.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	url    string // Discord webhook url
	client HttpClient
}

// NewClient create a Client with valid formatted webhook url.
func NewClient(hc HttpClient, url string) (*Client, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}
	return &Client{url: url, client: hc}, nil
}

// Send request to Discord webhook url via http post. Adjusted to the dynamic rate limit
func (c *Client) Send(messages []Message) ([]*http.Response, error) {
	dividedMessages := divideMessages(messages)
	if err := validateMessages(dividedMessages); err != nil {
		return nil, err
	}

	var responses []*http.Response
	for _, msg := range dividedMessages {
		resp, err := makeRequest(msg, c.url, c.client)
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
	contentType, body, err := writeBody(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return clt.Do(req)
}

// writeBody serialises Message. Returning error even though it's always nil to be consistent with multipartBody
func writeBody(msg Message) (contentType string, body io.Reader, err error) {
	if len(msg.Files) == 0 {
		// Marshal would never fail since Discord webhook message does not
		// contain types not supported by Marshal.
		payload, _ := json.Marshal(msg)
		return "application/json", bytes.NewBuffer(payload), nil
	}

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	if err = writePayload(msg, writer); err != nil {
		return
	}

	if err = writeFiles(msg.Files, writer); err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	return writer.FormDataContentType(), bytes.NewBuffer(b.Bytes()), nil
}

// writePayload writes Message to a part of the request body.
func writePayload(msg Message, writer *multipart.Writer) error {
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(h)
	if err != nil {
		return err
	}

	// Marshal would never fail since Discord webhook message does not
	// contain types not supported by Marshal.
	payload, _ := json.Marshal(msg)
	if _, err = part.Write(payload); err != nil {
		return err
	}
	return nil
}

// writeFiles writes each file to a part of the request body.
func writeFiles(files []*File, writer *multipart.Writer) error {
	escaper := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	for i, file := range files {
		h := textproto.MIMEHeader{}
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, escaper.Replace(file.Name)))
		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		h.Set("Content-Type", contentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			return err
		}

		if _, err = io.Copy(part, file.Reader); err != nil {
			return err
		}
	}
	return nil
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
