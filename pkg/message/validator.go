package message

import (
	"errors"
	"strings"
)

// These are all the length limit Discord API enforce on webhook message.
const (
	webhookEmbedsLimit = 10

	contentLimit = 2000

	embedTotalLimit       = 6000
	embedTitleLimit       = 256
	embedDescriptionLimit = 2048
	embedFieldsLimit      = 25
	embedAuthorNameLimit  = 256
	embedFooterTextLimit  = 2048

	fieldNameLimit  = 256
	fieldValueLimit = 1024
)

func limitError(field string) error {
	return errors.New(field + "length exceeding Discord API limit")
}

func validateField(f Field, embedLength *int) error {
	if len(f.Name) > fieldNameLimit {
		return limitError("Field name")
	}

	if len(f.Value) > fieldValueLimit {
		return limitError("Field value")
	}

	// Field name and value length is included in the embed total length
	*embedLength += len(f.Name)
	*embedLength += len(f.Value)

	if *embedLength > embedTotalLimit {
		return limitError("Embed total")
	}

	return nil
}

func validateEmbed(e Embed) error {
	var embedLength int

	if len(e.Title) > embedTitleLimit {
		return limitError("Embed title")
	}
	embedLength += len(e.Title)

	if len(e.Description) > embedDescriptionLimit {
		return limitError("Embed description")
	}
	embedLength += len(e.Description)

	if len(e.Author.Name) > embedAuthorNameLimit {
		return limitError("Embed author name")
	}
	embedLength += len(e.Author.Name)

	if len(e.Footer.Text) > embedFooterTextLimit {
		return limitError("Embed footer text")
	}
	embedLength += len(e.Footer.Text)

	if len(e.Fields) > embedFieldsLimit {
		return limitError("Embed field number")
	}

	for _, field := range e.Fields {
		err := validateField(field, &embedLength)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateURL(url string) error {
	webhookPrefix := "https://discord.com/api/webhooks/"
	if !strings.HasPrefix(url, webhookPrefix) {
		return errors.New("URL invalid")
	}

	return nil
}

// Validate checks Message object against Discord API limits. Returns slice
// containing length of each embed.
func Validate(m Message, url string) error {
	if m.Content == "" && len(m.Embeds) == 0 {
		return errors.New("Message must have either content or embeds")
	}

	if len(m.Content) > contentLimit {
		return limitError("Message content")
	}

	if len(m.Embeds) > webhookEmbedsLimit {
		return limitError("Message embed number")
	}

	if err := validateURL(url); err != nil {
		return err
	}

	for _, embed := range m.Embeds {
		err := validateEmbed(embed)
		if err != nil {
			return err
		}
	}

	return nil
}
