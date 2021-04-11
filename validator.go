package messenger

import (
	"errors"
	"strings"
)

// Limits Discord API enforces on webhook message.
const (
	// This library automatically divide message into smaller messages if it
	// exceeds the number of embed and embed total character limits.
	MessageEmbedNumLimit = 10
	MessageContentLimit  = 2000
	// Although not mentioned in documentation, in testings this limit seems
	// also enforced as the total character limit of multiple embeds in one
	// webhook message (excludes characters in "Content").
	EmbedTotalLimit       = 6000
	EmbedTitleLimit       = 256
	EmbedDescriptionLimit = 2048
	EmbedFieldNumLimit    = 25
	AuthorNameLimit       = 256
	FieldNameLimit        = 256
	FieldValueLimit       = 1024
	FooterTextLimit       = 2048
)

func limitError(field string) error {
	return errors.New(field + " length exceeding Discord API limit")
}

func validateFooter(f Footer) error {
	if f.Text == "" {
		return errors.New("Footer text is required")
	}

	if len(f.Text) > FooterTextLimit {
		return limitError("Embed footer text")
	}

	return nil
}

func validateField(f Field, embedLength *int) error {
	if f.Name == "" || f.Value == "" {
		return errors.New("Field name and value are required")
	}

	if len(f.Name) > FieldNameLimit {
		return limitError("Field name")
	}

	if len(f.Value) > FieldValueLimit {
		return limitError("Field value")
	}

	// Field name and value length is included in the embed total length
	*embedLength += len(f.Name)
	*embedLength += len(f.Value)

	if *embedLength > EmbedTotalLimit {
		return limitError("Embed total")
	}

	return nil
}

func validateEmbed(e Embed) error {
	var embedLength int

	if len(e.Title) > EmbedTitleLimit {
		return limitError("Embed title")
	}
	embedLength += len(e.Title)

	if len(e.Description) > EmbedDescriptionLimit {
		return limitError("Embed description")
	}
	embedLength += len(e.Description)

	if len(e.Author.Name) > AuthorNameLimit {
		return limitError("Embed author name")
	}
	embedLength += len(e.Author.Name)

	if e.Footer != (Footer{}) {
		err := validateFooter(e.Footer)
		if err != nil {
			return err
		}

		embedLength += len(e.Footer.Text)
	}

	if len(e.Fields) > EmbedFieldNumLimit {
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

// validateURL checks Discord webhook url validity.
func validateURL(url string) error {
	webhookPrefix := "https://discord.com/api/webhooks/"
	if !strings.HasPrefix(url, webhookPrefix) {
		return errors.New("URL invalid")
	}

	return nil
}

// validateMessage checks Message object against Discord API limits. Returns slice
// containing length of each embed.
func validateMessage(m Message) error {
	if m.Content == "" && len(m.Embeds) == 0 {
		return errors.New("Message must have either content or embeds")
	}

	if len(m.Content) > MessageContentLimit {
		return limitError("Message content")
	}

	if len(m.Embeds) > MessageEmbedNumLimit {
		return limitError("Message embed number")
	}

	for _, embed := range m.Embeds {
		err := validateEmbed(embed)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateRequest calls validateURL and validateMessage to check validity of a Request.
func validateRequest(r Request) (err error) {
	err = validateURL(r.URL)
	if err != nil {
		return
	}

	for _, msg := range r.Messages {
		err = validateMessage(msg)
		if err != nil {
			return
		}
	}

	return
}
