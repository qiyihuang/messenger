package messenger

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWait(t *testing.T) {
	t.Run("No limit", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-ratelimit-remaining", "")
		header.Set("x-ratelimit-reset-after", "")

		err := handleRateLimit(header)

		require.Equal(t, nil, err, "No limit failed")
	})

	t.Run("Atoi error", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-ratelimit-remaining", "")
		header.Set("x-ratelimit-reset-after", "something") // Avoid return by empty check

		err := handleRateLimit(header)

		_, ok := err.(*strconv.NumError)
		if !ok {
			t.Error("Atoi error failed")
		}
	})

	t.Run("Quota not exhausted", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-ratelimit-remaining", "1")
		header.Set("x-ratelimit-reset-after", "0")

		err := handleRateLimit(header)

		require.Equal(t, nil, err, "Quota not exhausted failed")
	})

	t.Run("ParseFloat error", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-ratelimit-remaining", "0")
		header.Set("x-ratelimit-reset-after", "")

		err := handleRateLimit(header)

		_, ok := err.(*strconv.NumError)
		if !ok {
			t.Error("ParseFloat error failed")
		}
	})

	t.Run("Return nil after sleep", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-ratelimit-remaining", "0")
		header.Set("x-ratelimit-reset-after", "0")

		err := handleRateLimit(header)

		require.Equal(t, nil, err, "Return nil after sleep failed")
	})
}
