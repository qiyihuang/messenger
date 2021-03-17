package messenger

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// sender mocks the Sender interface.
type sender struct{}

// send mocks the send method in Sender interface.
// Only returns resp, for the successful case.
func (s sender) send(p httpPoster) ([]*http.Response, error) {
	body, _ := json.Marshal("Ok")
	rr := httptest.NewRecorder()
	rr.Write(body)
	result := []*http.Response{rr.Result(), rr.Result()}
	return result, nil
}

func TestSend(t *testing.T) {
	t.Run("Return error", func(t *testing.T) {
		r := Request{Messages: []Message{{Content: "test"}}, URL: "wrong"}

		_, err := Send(r)

		require.Equal(t, errors.New("URL invalid"), err, "Return error failed")
	})

	t.Run("Return response", func(t *testing.T) {
		s := sender{}

		resp, err := Send(s)

		require.Equal(t, nil, err, "Return response no error failed")
		require.IsType(t, []*http.Response{}, resp, "Return response type failed")
	})
}
