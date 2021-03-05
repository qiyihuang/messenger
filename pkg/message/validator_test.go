package message

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errorMsg = "length exceeding Discord API limit"

func TestNewError(t *testing.T) {
	field := "test"

	err := limitError(field)

	require.Equal(t, errors.New(field+errorMsg), err, "New error failed")
}

func TestValidateField(t *testing.T) {
	t.Run("Field name limit", func(t *testing.T) {
		field := Field{Name: strings.Repeat("t", fieldNameLimit+1)}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Field name"+errorMsg), err, "Field name limit failed")
	})

	t.Run("Field value limit", func(t *testing.T) {
		field := Field{Value: strings.Repeat("t", fieldValueLimit+1)}
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
		field := Field{}
		length := embedTotalLimit + 1

		err := validateField(field, &length)

		require.Equal(t, errors.New("Embed total"+errorMsg), err, "Embed total limit failed")
	})

	t.Run("No error", func(t *testing.T) {
		field := Field{Name: "Pass", Value: "Pass"}
		length := 1

		err := validateField(field, &length)

		require.Equal(t, nil, err, "No error failed")
	})
}
