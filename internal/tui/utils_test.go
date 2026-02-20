package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeftShiftString(t *testing.T) {
	emptyString := leftShiftString("")
	assert.Equal(t, "", emptyString, "expected emptyString")
}
