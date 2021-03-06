package message

import (
	"context"
	"errors"
	"testing"

	"github.com/AnilRedshift/captions_please_go/pkg/structured_error"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLoadMessages(t *testing.T) {
	assert.NoError(t, LoadMessages())
}

func TestGetErrorMessage(t *testing.T) {
	assert.NoError(t, LoadMessages())
	anError := errors.New("oh no")
	tests := []struct {
		name     string
		err      structured_error.StructuredError
		enResult Localized
	}{
		{
			name:     "Defaults to an unknown error",
			err:      structured_error.Wrap(anError, structured_error.ErrorType(999)),
			enResult: unknownErrorFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.enResult, ErrorMessage(context.Background(), test.err))
		})
	}
}

func TestGerman(t *testing.T) {
	assert.NoError(t, LoadMessages())
	ctx := WithLanguage(context.Background(), language.German)
	assert.Equal(t, Localized("Hilfe"), sprint(ctx, "help"))
}

func TestLabelImage(t *testing.T) {
	assert.NoError(t, LoadMessages())
	assert.Equal(t, Localized("Image 1: foo"), LabelImage(context.Background(), Unlocalized("foo"), 0))
	assert.Equal(t, Localized("Image 2: foo"), LabelImage(context.Background(), Unlocalized("foo"), 1))
}
