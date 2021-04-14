package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountEmbed(t *testing.T) {
	var total int16 = 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8
	embed := Embed{
		Title:       strings.Repeat("t", 1),
		Description: strings.Repeat("t", 2),
		Author:      Author{Name: strings.Repeat("t", 3)},
		Footer:      Footer{Text: strings.Repeat("t", 4)},
		Fields: []Field{
			{Name: strings.Repeat("t", 5), Value: strings.Repeat("t", 6)},
			{Name: strings.Repeat("t", 7), Value: strings.Repeat("t", 8)},
		},
	}

	count := countEmbed(embed)

	require.Equal(t, total, count, "CountEmbed failed")
}

func TestDivideEmbeds(t *testing.T) {
	t.Run("Divide by embed character limit", func(t *testing.T) {
		expectedNumber := 3 //1000 + 2000 + 3000, 4000, 4000
		embeds := []Embed{
			{Description: strings.Repeat("t", 1000)},
			{Description: strings.Repeat("e", 2000)},
			{Description: strings.Repeat("s", 3000)},
			{Description: strings.Repeat("t", 4000)},
			{Description: strings.Repeat("t", 4000)},
		}
		msg := Message{Username: "t", Content: "test", Embeds: embeds}

		dividedEmbeds := divideEmbeds(msg)

		require.Equal(t, expectedNumber, len(dividedEmbeds), "Divide by embed character limit failed")
	})

	t.Run("Divide by embed number", func(t *testing.T) {
		expectedNumber := 2
		embeds := []Embed{
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
			{Description: "t"},
		}
		msg := Message{Username: "t", Content: "test", Embeds: embeds}

		dividedEmbeds := divideEmbeds(msg)

		require.Equal(t, expectedNumber, len(dividedEmbeds), "Divide by embed number failed")
	})
}

func TestDivideMessages(t *testing.T) {
	t.Run("Content in only 1 message", func(t *testing.T) {
		content := "test"
		embeds := []Embed{
			{Description: strings.Repeat("t", 4000)},
			{Description: strings.Repeat("e", 4000)},
			{Description: strings.Repeat("s", 4000)},
		}
		msgs := []Message{
			{Username: "t", Content: content, Embeds: embeds},
		}

		dividedMsgs := divideMessages(msgs)

		require.Equal(t, content, dividedMsgs[0].Content, "Content in first message failed")
		require.Equal(t, "", dividedMsgs[1].Content, "Content in second message failed")
		require.Equal(t, "", dividedMsgs[2].Content, "Content in third message failed")
	})
}

func TestFormatBody(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		msg := Message{}
		jsonMsg, _ := json.Marshal(msg)
		expectedBody := bytes.NewBuffer(jsonMsg)

		body := formatBody(msg)

		require.Equal(t, expectedBody, body, "Json marshal no error failed")
	})
}

func TestRespError(t *testing.T) {
	type body struct {
		Message string `json:"message,omitempty"`
		Other   string `json:"other,omitempty"`
	}

	t.Run("Has error in response", func(t *testing.T) {
		rawBody := body{Message: "test"}
		jsonBody, _ := json.Marshal(rawBody)
		rr := httptest.NewRecorder()
		rr.Write(jsonBody)

		err := respError(rr.Result())

		require.Equal(t, errors.New("Discord API error: test"), err, "Has error in response failed")
	})

	t.Run("Decode return if EOF", func(t *testing.T) {
		rr := httptest.NewRecorder()
		// Write empty body to trigger decode EOF error.
		rr.Write(nil)

		err := respError(rr.Result())

		require.IsType(t, nil, err, "Decode return if EOF failed")
	})

	t.Run("Decode error", func(t *testing.T) {
		body, _ := json.Marshal(1)
		rr := httptest.NewRecorder()
		// Write empty body to trigger decode EOF error.
		rr.Write(body)

		err := respError(rr.Result())

		require.IsType(t, &json.UnmarshalTypeError{}, err, "Decode error failed")
	})

	t.Run("No error", func(t *testing.T) {
		rawBody := body{Other: "Ok"}
		jsonBody, _ := json.Marshal(rawBody)
		rr := httptest.NewRecorder()
		rr.Write(jsonBody)

		err := respError(rr.Result())

		require.Equal(t, nil, err, "Resp no error failed")
	})
}

// postError imitates error returned by http.Post
func postError(url string, contentType string, body io.Reader) (*http.Response, error) {
	return nil, errors.New("Post error")
}

// responseError imitates response with Discord error.
func responseError(url string, contentType string, body io.Reader) (*http.Response, error) {
	respBody := struct {
		Message string `json:"message,omitempty"`
	}{"Response error"}
	jsonBody, _ := json.Marshal(respBody)
	rr := httptest.NewRecorder()
	rr.Write(jsonBody)
	return rr.Result(), nil
}

// waitError imitates error thrown by ratelimit.Wait.
func waitError(url string, contentType string, body io.Reader) (*http.Response, error) {
	rr := httptest.NewRecorder()
	header := rr.Header()
	header.Add("x-ratelimit-remaining", "wrong")
	header.Add("x-ratelimit-reset-after", "1")
	respBody, _ := json.Marshal(struct{ Other string }{"Ok"})
	rr.Write(respBody)
	return rr.Result(), nil

}

// noError imitates successful response from http.Post.
func noError(url string, contentType string, body io.Reader) (*http.Response, error) {
	jsonBody, _ := json.Marshal(struct{ Message string }{"Ok"})
	rr := httptest.NewRecorder()
	rr.Write(jsonBody)
	return rr.Result(), nil
}

func TestRequestSend(t *testing.T) {
	t.Run("validateURL error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "test"}}, URL: "wrong"}

		_, err := r.Send()

		require.Equal(t, errors.New("URL invalid"), err, "validateURL error failed")
	})

	t.Run("validateMessage error", func(t *testing.T) {
		r := Request{Messages: []Message{{}}, URL: "https://discord.com/api/webhooks/"}

		_, err := r.Send()

		require.Equal(t, errors.New("Message must have either content or embeds"), err, "validateMessage error failed")
	})

	t.Run("Post error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		post = postError

		_, err := r.Send()

		require.Equal(t, errors.New("Post error"), err, "Post error failed")
	})

	t.Run("respError error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		post = responseError

		_, err := r.Send()

		require.Equal(t, errors.New("Discord API error: Response error"), err, "respError error failed")
	})

	t.Run("ratelimit.Wait error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		post = waitError

		_, err := r.Send()

		require.IsType(t, &strconv.NumError{}, err, "respError error failed")
	})

	t.Run("Success", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		post = noError

		_, err := r.Send()

		require.Equal(t, nil, err, "Success failed")
	})
}
