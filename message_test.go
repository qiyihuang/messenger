package messenger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
