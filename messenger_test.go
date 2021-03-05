package messenger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Run("validateMessage error", func(t *testing.T) {
		url := "https://discord.com/api/webhooks/"
		msg := Message{}

		_, err := Send(url, msg)

		require.Equal(t, errors.New("Message must have either content or embeds"), err, "validateMessage error failed")
	})

	t.Run("makeRequest error", func(t *testing.T) {
		url := "wrong"
		msg := Message{Content: "test"}

		_, err := Send(url, msg)

		require.Equal(t, errors.New("URL invalid"), err, "Validate error failed")
	})
}
