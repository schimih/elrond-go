package manager

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatterFactory(t *testing.T) {
	fF := NewFormatterFactory()

	assert.False(t, check.IfNil(fF))
}

func TestCreateFormatter_TableFormatter(t *testing.T) {
	fF := NewFormatterFactory()
	TableFormatter := &TableFormatter{
		header:    make([]string, 0),
		dataLines: make([]*display.LineData, 0)}

	formatter := fF.CreateFormatter(utils.Table)
	require.Equal(t, TableFormatter, formatter)
}

func TestCreateFormatter_JsonFormatter(t *testing.T) {
	fF := NewFormatterFactory()
	TableFormatter := &JsonFormatter{}

	formatter := fF.CreateFormatter(utils.JSON)
	require.Equal(t, TableFormatter, formatter)
}

func TestCreateFormatter_XMLFormatter(t *testing.T) {
	fF := NewFormatterFactory()
	TableFormatter := &XMLFormatter{}

	formatter := fF.CreateFormatter(utils.XML)
	require.Equal(t, TableFormatter, formatter)
}

func TestCreateFormatter_GINFormatter(t *testing.T) {
	fF := NewFormatterFactory()

	formatter := fF.CreateFormatter(utils.GIN)
	require.Nil(t, formatter)
}
