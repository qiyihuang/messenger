package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		url := "wrong"

		c, err := NewClient(http.DefaultClient, url)

		require.Equal(t, (*Client)(nil), c, "TestNewClient error failed")
		require.EqualError(t, err, "invalid webhook URL")
	})

	t.Run("success", func(t *testing.T) {
		url := "https://discord.com/api/webhooks/something"

		_, err := NewClient(http.DefaultClient, url)

		require.NoError(t, err)
	})
}

func TestClientSend(t *testing.T) {
	t.Run("validateMessages error", func(t *testing.T) {
		// %% will fail makeRequest
		c := &Client{url: "ok", client: http.DefaultClient}

		_, err := c.Send([]Message{})

		require.Error(t, err)
	})

	t.Run("makeRequest error", func(t *testing.T) {
		// %% will fail makeRequest
		c := &Client{url: "%%", client: http.DefaultClient}

		_, err := c.Send([]Message{{Content: "Ok"}})

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

		c := &Client{url: server.URL, client: http.DefaultClient}

		_, err := c.Send([]Message{{Content: "Ok"}})

		require.Equal(t, errors.New("Discord API error: test error"), err, "respError error failed")
	})

	t.Run("ratelimit.Wait error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("x-ratelimit-remaining", "a") // trigger strconv error.
		}))
		defer server.Close()

		c := &Client{url: server.URL, client: http.DefaultClient}

		_, err := c.Send([]Message{{Content: "Ok"}})

		require.IsType(t, &strconv.NumError{}, err, "respError error failed")
	})

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		c := &Client{url: server.URL, client: http.DefaultClient}

		_, err := c.Send([]Message{{Content: "Ok"}})

		require.NoError(t, err)
	})
}

func TestMakeRequest(t *testing.T) {
	t.Run("multipartBody no error", func(t *testing.T) {
		msg := Message{Files: []*File{
			{Name: "Test", Reader: bytes.NewBuffer([]byte{1})},
		}}
		clt := &http.Client{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		_, err := makeRequest(msg, server.URL, clt)

		require.NoError(t, err, "multipartBody no error failed")
	})

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

func TestWriteBody(t *testing.T) {
	t.Run("Normal body", func(t *testing.T) {
		msg := Message{Content: "test"}

		contentType, body, err := writeBody(msg)

		require.Equal(t, contentType, "application/json", "Normal body failed")
		require.NotEmpty(t, body, "Normal body failed")
		require.NoError(t, err, "Normal body failed")
	})

	t.Run("Multipart body", func(t *testing.T) {
		msg := Message{
			Content: "test",
			Files: []*File{
				{Name: "test.jpg", Reader: bytes.NewBuffer([]byte{1})},
			},
		}

		contentType, body, err := writeBody(msg)

		require.True(t, strings.Contains(contentType, "multipart/form-data"), "Multipart body failed")
		require.NotEmpty(t, body, "Multipart body failed")
		require.NoError(t, err, "Multipart body failed")
	})
}

type writerMock struct {
}

func (w writerMock) Write(p []byte) (n int, err error) {
	return 0, errors.New("Test")
}

func TestWritePayload(t *testing.T) {
	t.Run("CreatePart error", func(t *testing.T) {
		msg := Message{Content: "something"}
		w := writerMock{}
		mpw := multipart.NewWriter(w)

		err := writePayload(msg, mpw)

		require.Error(t, err, "Test", "CreatePart error failed")
	})

	t.Run("No error", func(t *testing.T) {
		msg := Message{Content: "something"}
		mpw := multipart.NewWriter(&bytes.Buffer{})

		err := writePayload(msg, mpw)

		require.NoError(t, err, "No error failed")
	})
}

func TestWriteFiles(t *testing.T) {
	t.Run("CreatePart error", func(t *testing.T) {
		files := []*File{
			{
				Name: "test.jpg",
			},
		}
		w := writerMock{}
		mpw := multipart.NewWriter(w)

		err := writeFiles(files, mpw)

		require.Error(t, err, "Test", "CreatePart error failed")
	})

	t.Run("No error", func(t *testing.T) {
		// No io.Reader to copy from.
		files := []*File{
			{
				Name:   "test.jpg",
				Reader: bytes.NewBuffer([]byte{1}),
			},
			{
				Name:        "test2.jpg",
				ContentType: "image/jpeg",
				Reader:      bytes.NewBuffer([]byte{1}),
			},
		}
		mpw := multipart.NewWriter(&bytes.Buffer{})

		err := writeFiles(files, mpw)

		require.NoError(t, err, "No error failed")
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
