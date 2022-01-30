package scan

import (
	"testing"

	"github.com/elrond-go/cmd/vat/utils"
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
		Grammar:           utils.TCP_ELROND,
		SlicedParsedInput: make([]*gonmap.NmapRun, 0),
	}

	parsingResults := parserData.Parse()

	assert.NotNil(t, parsingResults)
}
