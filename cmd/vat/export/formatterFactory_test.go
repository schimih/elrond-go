package export

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

	formatter, _ := fF.CreateFormatter(utils.Table)
	require.Equal(t, TableFormatter, formatter)
}

func TestCreateFormatter_JsonFormatter(t *testing.T) {
	fF := NewFormatterFactory()
	Jsonformatter := &JsonFormatter{}

	formatter, _ := fF.CreateFormatter(utils.JSON)
	require.Equal(t, Jsonformatter, formatter)
}

func TestCreateFormatter_XMLFormatter(t *testing.T) {
	fF := NewFormatterFactory()
	XMLFormatter := &XMLFormatter{}

	formatter, _ := fF.CreateFormatter(utils.XML)
	require.Equal(t, XMLFormatter, formatter)
}

func TestCreateFormatter_GINFormatter(t *testing.T) {
	fF := NewFormatterFactory()

	formatter, _ := fF.CreateFormatter(utils.GIN)
	require.Nil(t, formatter)
}
