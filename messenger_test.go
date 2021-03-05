package messenger

import (
	"errors"
	"testing"

	"github.com/qiyihuang/messenger/pkg/message"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Run("Validate error", func(t *testing.T) {
		url := "wrong"
		msg := message.Message{Content: "test"}

		err := Send(url, msg)

		require.Equal(t, errors.New("URL invalid"), err, "Validate error failed")
	})
}
