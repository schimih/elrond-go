package factory

import (
	"testing"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateParser(t *testing.T) {
	newParserFactory := NewParserFactory()
	input := make([][]byte, 1)
	parser := newParserFactory.CreateParser(input, utils.TCP_ELROND)

	assert.NotNil(t, parser)
	assert.IsType(t, &scan.ParserData{}, parser)
}
