package factory

import (
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type ParserFactory struct {
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{}
}

func (pf *ParserFactory) CreateParser(input [][]byte, grammar utils.AnalysisType) scan.Parser {
	return &scan.ParserData{
		Input:             input,
		AnalyzedTargets:   make([]scan.ScannedTarget, 0),
		Grammar:           grammar,
		SlicedParsedInput: make([]*go_nmap.NmapRun, 0),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *ParserFactory) IsInterfaceNil() bool {
	return pf == nil
}
