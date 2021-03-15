package messenger

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type sender struct{}

var fakeResp *http.Response

// Only returns resp, for the successful case.
func (s sender) send() (*http.Response, error) {
	body, _ := json.Marshal("test")
	rr := httptest.NewRecorder()
	rr.Write(body)
	fakeResp = rr.Result()
	return fakeResp, nil
}

func TestSend(t *testing.T) {
	t.Run("Return error", func(t *testing.T) {
		r := Request{Msg: Message{Content: "test"}, URL: "wrong"}

		_, err := Send(r)

		require.Equal(t, errors.New("URL invalid"), err, "Return error failed")
	})

	t.Run("Return response", func(t *testing.T) {
		s := sender{}

		resp, _ := Send(s)

		require.Equal(t, fakeResp, resp, "Return response failed")
	})
}
