package messenger

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errorMsg = " length exceeding Discord API limit"

func TestNewError(t *testing.T) {
	field := "test"

	err := limitError(field)

	require.Equal(t, errors.New(field+errorMsg), err, "New error failed")
}

func TestValidateFooter(t *testing.T) {
	t.Run("Text empty", func(t *testing.T) {
		footer := Footer{IconURL: "url"}

		err := validateFooter(footer)

		require.Equal(t, errors.New("Footer text is required"), err, "Text empty failed")
	})

	t.Run("Embed footer limit", func(t *testing.T) {
		footer := Footer{Text: strings.Repeat("t", FooterTextLimit+1)}

		err := validateFooter(footer)

		require.Equal(t, errors.New("Embed footer text"+errorMsg), err, "Embed footer limit failed")
	})

	t.Run("Pass", func(t *testing.T) {
		footer := Footer{Text: "Ok"}

		err := validateFooter(footer)

		require.Equal(t, nil, err, "Pass failed")
	})
}

func TestValidateField(t *testing.T) {
	t.Run("Name empty", func(t *testing.T) {
		field := Field{Value: "Ok"}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Field name and value are required"), err, "Name empty failed")
	})

	t.Run("Value empty", func(t *testing.T) {
		field := Field{Name: "Ok"}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Field name and value are required"), err, "Name empty failed")
	})

	t.Run("Field name limit", func(t *testing.T) {
		field := Field{Name: strings.Repeat("t", FieldNameLimit+1), Value: "Ok"}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Field name"+errorMsg), err, "Field name limit failed")
	})

	t.Run("Field value limit", func(t *testing.T) {
		field := Field{Name: "Ok", Value: strings.Repeat("t", FieldValueLimit+1)}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Field value"+errorMsg), err, "Field value limit failed")
	})

	t.Run("embedLength addition", func(t *testing.T) {
		originalLength := 1
		nameLength := 2
		valueLength := 3
		field := Field{Name: strings.Repeat("t", nameLength), Value: strings.Repeat("t", valueLength)}
		length := 1

		validateField(field, &length)

		require.Equal(t, originalLength+nameLength+valueLength, length, "embedLength addition failed")
	})

	t.Run("Embed total limit", func(t *testing.T) {
		field := Field{Name: "Ok", Value: "Ok"}
		length := EmbedTotalLimit + 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Embed total"+errorMsg), err, "Embed total limit failed")
	})

	t.Run("No error", func(t *testing.T) {
		field := Field{Name: "Ok", Value: "Ok"}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, nil, err, "No error failed")
	})
}

func TestValidateEmbed(t *testing.T) {
	t.Run("Embed title limit", func(t *testing.T) {
		embed := Embed{Title: strings.Repeat("t", EmbedTitleLimit+1)}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Embed title"+errorMsg), err, "Embed title limit failed")
	})

	t.Run("Embed description limit", func(t *testing.T) {
		embed := Embed{Description: strings.Repeat("t", EmbedDescriptionLimit+1)}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Embed description"+errorMsg), err, "Embed description limit failed")
	})

	t.Run("Embed author name limit", func(t *testing.T) {
		embed := Embed{Author: Author{Name: strings.Repeat("t", AuthorNameLimit+1)}}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Embed author name"+errorMsg), err, "Embed author name limit failed")
	})

	t.Run("Validate footer", func(t *testing.T) {
		embed := Embed{Footer: Footer{IconURL: "no text should fail"}}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Footer text is required"), err, "Validate footer failed")
	})

	t.Run("Validate fields number", func(t *testing.T) {
		// 26 fields
		fields := []Field{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
		embed := Embed{Fields: fields}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Embed field number"+errorMsg), err, "Validate fields number failed")
	})

	t.Run("Validate fields", func(t *testing.T) {
		fields := []Field{{Name: "Ok", Value: "Ok"}, {Name: "Ok"}}
		embed := Embed{Fields: fields}

		err := validateEmbed(embed)

		require.Equal(t, errors.New("Field name and value are required"), err, "Validate fields failed")
	})

	t.Run("Pass", func(t *testing.T) {
		embed := Embed{}

		err := validateEmbed(embed)

		require.Equal(t, nil, err, "Pass failed")
	})

	t.Run("embedLength addition", func(t *testing.T) {
		passedEmbed := Embed{
			Title:       "t",
			Description: "t",
			Author:      Author{Name: "t"},
			Footer:      Footer{Text: "t"},
			Fields: []Field{
				{Name: "t", Value: strings.Repeat("t", 1000)},
				{Name: "t", Value: strings.Repeat("t", 1000)},
				{Name: "t", Value: strings.Repeat("t", 1000)},
				{Name: "t", Value: strings.Repeat("t", 1000)},
				{Name: "t", Value: strings.Repeat("t", 1000)},
				{Name: "t", Value: strings.Repeat("t", 981)},
				// The total length 5991 makes it for every single "t"
				// field +1, the embed fails by just 1 char, so every field
				// must be counted in order to pass
			},
		}
		failedEmbed := Embed{
			Title:       "tt",
			Description: "tt",
			Author:      Author{Name: "tt"},
			Footer:      Footer{Text: "tt"},
			Fields: []Field{
				{Name: "tt", Value: strings.Repeat("t", 1000)},
				{Name: "tt", Value: strings.Repeat("t", 1000)},
				{Name: "tt", Value: strings.Repeat("t", 1000)},
				{Name: "tt", Value: strings.Repeat("t", 1000)},
				{Name: "tt", Value: strings.Repeat("t", 1000)},
				{Name: "tt", Value: strings.Repeat("t", 981)},
			},
		}

		err := validateEmbed(passedEmbed)
		require.Equal(t, nil, err, "embedLength addition pass failed")

		err = validateEmbed(failedEmbed)
		require.Equal(t, errors.New("Embed total"+errorMsg), err, "embedLength addition error failed")
	})
}

func TestValidateURL(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		url := "wrong"

		err := validateURL(url)

		require.Equal(t, errors.New("invalid webhook URL"), err, "Invalid failed")
	})

	t.Run("Pass", func(t *testing.T) {
		url := "https://discord.com/api/webhooks/"

		err := validateURL(url)

		require.Equal(t, nil, err, "Pass failed")
	})
}

func TestValidateMessage(t *testing.T) {
	t.Run("Neither content nor embeds", func(t *testing.T) {
		msg := Message{}

		err := validateMessage(msg)

		require.Equal(t, errors.New("Message must have either content or embeds"), err, "Neither content nor embeds failed")
	})

	t.Run("Content limit", func(t *testing.T) {
		msg := Message{Content: strings.Repeat("t", MessageContentLimit+1)}

		err := validateMessage(msg)

		require.Equal(t, errors.New("Message content"+errorMsg), err, "Content limit failed")
	})

	t.Run("Embed number limit", func(t *testing.T) {
		embeds := []Embed{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
		msg := Message{Embeds: embeds}

		err := validateMessage(msg)

		require.Equal(t, errors.New("Message embed number"+errorMsg), err, "Embed number limit failed")
	})

	t.Run("Validate embeds", func(t *testing.T) {
		embeds := []Embed{{}, {Title: strings.Repeat("t", EmbedTitleLimit+1)}}
		msg := Message{Embeds: embeds}

		err := validateMessage(msg)

		require.Equal(t, errors.New("Embed title"+errorMsg), err, "Validate embeds failed")
	})

	t.Run("Pass", func(t *testing.T) {
		msg := Message{Content: "Ok"}

		err := validateMessage(msg)

		require.Equal(t, nil, err, "Pass failed")
	})
}

func TestValidateMessages(t *testing.T) {
	t.Run("No message", func(t *testing.T) {
		msgs := []Message{}

		err := validateMessages(msgs)

		require.EqualError(t, err, "request must have a least 1 message")

	})

	t.Run("Message error", func(t *testing.T) {
		msgs := []Message{{Content: "Ok"}, {}} // Failed on second make sure it loops.

		err := validateMessages(msgs)

		require.EqualError(t, err, "Message must have either content or embeds")
	})

	t.Run("Pass", func(t *testing.T) {
		msgs := []Message{{Content: "Ok"}}

		err := validateMessages(msgs)

		require.Equal(t, nil, err, "Message error failed")
	})
}
