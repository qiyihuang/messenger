package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
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

	t.Run("Decode error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		// Write empty body to trigger decode EOF error.
		rr.Write(nil)

		err := respError(rr.Result())

		require.Equal(t, errors.New("EOF"), err, "Decode error failed")
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

func TestRequestSend(t *testing.T) {
	t.Run("validateURL error", func(t *testing.T) {
		r := Request{Msg: Message{Content: "test"}, URL: "wrong"}

		_, err := r.send()

		require.Equal(t, errors.New("URL invalid"), err, "validateURL error failed")
	})

	t.Run("validateMessage error", func(t *testing.T) {
		r := Request{Msg: Message{}, URL: "https://discord.com/api/webhooks/"}

		_, err := r.send()

		require.Equal(t, errors.New("Message must have either content or embeds"), err, "validateMessage error failed")
	})
}
