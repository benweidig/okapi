package okapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Notifier_Set(t *testing.T) {
	// ASSERT
	assert.Nil(t, notifier)

	// ARRANGE
	ch := make(chan Info)

	// ACT
	Notifier(ch)

	// ASSERT
	assert.NotNil(t, notifier)
}
