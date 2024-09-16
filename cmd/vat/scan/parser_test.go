package scan

import (
	"testing"

	core "github.com/elrond-go/cmd/vat/core"
	gonmap "github.com/lair-framework/go-nmap"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	text := "test"
	textAsByte := []byte(text)
	input := [][]byte{textAsByte}
	parserData := &ParserData{
		Input:             input,
		AnalyzedTargets:   make([]ScannedTarget, 0),
		Grammar:           core.TCP_ELROND,
		SlicedParsedInput: make([]*gonmap.NmapRun, 0),
	}

	parsingResults := parserData.Parse()

	assert.NotNil(t, parsingResults)
}
