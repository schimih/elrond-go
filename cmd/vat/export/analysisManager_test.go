package export

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalysisManager(t *testing.T) {

	aM, err := NewAnalysisManager()

	assert.False(t, check.IfNil(aM))
	assert.Nil(t, err)
}

func TestNewAnalysisManager_FormatterNilCheck(t *testing.T) {
	//fakeFF := &FormatterFactory{}

	aM, err := NewAnalysisManager()

	assert.True(t, check.IfNil(aM))
	expectedErrorString := "FormatterFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}
