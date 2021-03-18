package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

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

// mockedPoster mocks the HttpPoster interface.
type mockedPoster struct {
	// flags to control response.
	postError bool
	respError bool
	waitError bool
}

// Post mocks Post method in HttpPoster interface.
func (mp mockedPoster) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	// Imitate http.Post error.
	if mp.postError == true {
		return nil, errors.New("Post error")
	}

	// Imitate response with Discord error.
	if mp.respError == true {
		respBody := struct {
			Message string `json:"message,omitempty"`
		}{"Response error"}
		jsonBody, _ := json.Marshal(respBody)
		rr := httptest.NewRecorder()
		rr.Write(jsonBody)
		return rr.Result(), nil
	}

	if mp.waitError == true {
		rr := httptest.NewRecorder()
		header := rr.Header()
		header.Add("x-ratelimit-remaining", "wrong")
		header.Add("x-ratelimit-reset-after", "1")
		body, _ := json.Marshal(struct{ Other string }{"Ok"})
		rr.Write(body)
		return rr.Result(), nil
	}

	jsonBody, _ := json.Marshal(struct{ Message string }{"Ok"})
	rr := httptest.NewRecorder()
	rr.Write(jsonBody)
	return rr.Result(), nil
}

func TestRequestSend(t *testing.T) {
	t.Run("validateURL error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "test"}}, URL: "wrong"}

		_, err := r.send(http.DefaultClient)

		require.Equal(t, errors.New("URL invalid"), err, "validateURL error failed")
	})

	t.Run("validateMessage error", func(t *testing.T) {
		r := Request{Messages: []Message{{}}, URL: "https://discord.com/api/webhooks/"}

		_, err := r.send(http.DefaultClient)

		require.Equal(t, errors.New("Message must have either content or embeds"), err, "validateMessage error failed")
	})

	t.Run("Post error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		mp := mockedPoster{postError: true, respError: false, waitError: false}

		_, err := r.send(mp)

		require.Equal(t, errors.New("Post error"), err, "Post error failed")
	})

	t.Run("respError error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		mp := mockedPoster{postError: false, respError: true, waitError: false}

		_, err := r.send(mp)

		require.Equal(t, errors.New("Discord API error: Response error"), err, "respError error failed")
	})

	t.Run("ratelimit.Wait error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		mp := mockedPoster{postError: false, respError: false, waitError: true}

		_, err := r.send(mp)

		require.IsType(t, &strconv.NumError{}, err, "respError error failed")
	})

	t.Run("Success", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "Ok"}}, URL: "https://discord.com/api/webhooks/"}
		mp := mockedPoster{postError: false, respError: false, waitError: false}

		_, err := r.send(mp)

		require.Equal(t, nil, err, "Success failed")
	})
}
