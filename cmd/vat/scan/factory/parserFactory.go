package factory

import (
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type ParserFactory struct {
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{}
}

func (pf *ParserFactory) CreateParser(input [][]byte, grammar utils.AnalysisType) scan.Parser {
	return &scan.ParserData{
		Input:           input,
		AnalyzedTargets: make([]scan.ScannedTarget, 0),
		Grammar:         grammar,
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *ParserFactory) IsInterfaceNil() bool {
	return pf == nil
}
