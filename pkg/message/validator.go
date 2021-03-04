package message

import "errors"

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

func newError(field string) error {
	return errors.New(field + "length exceeding Discord API limit")
}

func validateField(f Field, embedLength *int) error {
	if len(f.Name) > fieldNameLimit {
		return newError("Field name")
	}

	if len(f.Value) > fieldValueLimit {
		return newError("Field value")
	}

	// Field name and value length is included in the embed total length
	*embedLength += len(f.Name)
	*embedLength += len(f.Value)

	if *embedLength > embedTotalLimit {
		return newError("Embed total")
	}

	return nil
}

func validateEmbed(e Embed) error {
	var embedLength int

	if len(e.Title) > embedTitleLimit {
		return newError("Embed title")
	}
	embedLength += len(e.Title)

	if len(e.Description) > embedDescriptionLimit {
		return newError("Embed description")
	}
	embedLength += len(e.Description)

	if len(e.Author.Name) > embedAuthorNameLimit {
		return newError("Embed author name")
	}
	embedLength += len(e.Author.Name)

	if len(e.Footer.Text) > embedFooterTextLimit {
		return newError("Embed footer text")
	}
	embedLength += len(e.Footer.Text)

	if len(e.Fields) > embedFieldsLimit {
		return newError("Embed field number")
	}

	for _, field := range e.Fields {
		err := validateField(field, &embedLength)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate checks Message object against Discord API limits. Returns slice
// containing length of each embed.
func Validate(m Message) error {
	if m.Content == "" && len(m.Embeds) == 0 {
		return errors.New("Message must have either content or embeds")
	}

	if len(m.Content) > contentLimit {
		return newError("Message content")
	}

	if len(m.Embeds) > webhookEmbedsLimit {
		return newError("Message embed number")
	}

	for _, embed := range m.Embeds {
		err := validateEmbed(embed)
		if err != nil {
			return err
		}
	}

	return nil
}
