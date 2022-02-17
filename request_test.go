package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		msgs := []Message{}
		url := "https://discord.com/api/webhooks/something"

		req, err := NewRequest(http.DefaultClient, url, msgs)

		require.Equal(t, (*request)(nil), req, "TestNewRequest error failed")
		require.EqualError(t, err, "request must have a least 1 message")
	})

	t.Run("success", func(t *testing.T) {
		msgs := []Message{{Content: "test"}}
		url := "https://discord.com/api/webhooks/something"

		_, err := NewRequest(http.DefaultClient, url, msgs)

		require.NoError(t, err)
	})
}

func TestRequestSend(t *testing.T) {
	t.Run("makeRequest error", func(t *testing.T) {
		// %% will fail makeRequest
		r := request{messages: []Message{{Content: "Ok"}}, url: "%%", client: http.DefaultClient}

		_, err := r.Send()

		require.Error(t, err)
	})

	t.Run("respError error", func(t *testing.T) {
		// Return a payload containing error message
		resp := make(map[string]string)
		resp["message"] = "test error"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(resp)
			w.Write(j)
		}))
		defer server.Close()

		r := request{messages: []Message{{Content: "Ok"}}, url: server.URL, client: http.DefaultClient}

		_, err := r.Send()

		require.Equal(t, errors.New("Discord API error: test error"), err, "respError error failed")
	})

	t.Run("ratelimit.Wait error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("x-ratelimit-remaining", "a") // trigger strconv error.
		}))
		defer server.Close()

		r := request{messages: []Message{{Content: "Ok"}}, url: server.URL, client: http.DefaultClient}

		_, err := r.Send()

		require.IsType(t, &strconv.NumError{}, err, "respError error failed")
	})

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		r := request{messages: []Message{{Content: "Ok"}}, url: server.URL, client: http.DefaultClient}

		_, err := r.Send()

		require.NoError(t, err)
	})
}

func TestMakeRequest(t *testing.T) {
	t.Run("NewRequest error", func(t *testing.T) {
		msg := Message{}
		url := "%%" // This will make NewRequest failed
		clt := &http.Client{}

		resp, err := makeRequest(msg, url, clt)

		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("Do error", func(t *testing.T) {
		msg := Message{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close() // Close before req sent
		clt := server.Client()

		_, err := makeRequest(msg, server.URL, clt)

		require.Error(t, err)
	})

	t.Run("No error", func(t *testing.T) {
		msg := Message{Content: "hi"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()
		clt := server.Client()

		_, err := makeRequest(msg, server.URL, clt)

		require.NoError(t, err)
	})
}

func TestCountEmbed(t *testing.T) {
	var total = 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8
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
		expectedNumber := 3 //1000 + 2000 + 3000, 3000, 4000 + 2000
		embeds := []Embed{
			{Description: strings.Repeat("t", 1000)},
			{Description: strings.Repeat("e", 2000)},
			{Description: strings.Repeat("s", 3000)},
			{Description: strings.Repeat("t", 3000)},
			{Description: strings.Repeat("t", 4000)},
			{Description: strings.Repeat("t", 2000)},
		}
		msg := Message{Username: "t", Content: "test", Embeds: embeds}

		dividedEmbeds := divideEmbeds(msg)

		require.Equal(t, expectedNumber, len(dividedEmbeds), "Divide by embed character limit failed")
	})

	t.Run("Divide by embed number", func(t *testing.T) {
		expectedNumber := 3
		embeds := []Embed{ // 21 embeds
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
