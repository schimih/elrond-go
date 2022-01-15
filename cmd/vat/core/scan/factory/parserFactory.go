package factory

import (
	"github.com/elrond-go/cmd/vat/core/scan"
)

type ParserFactory struct {
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{}
}

func (pf *ParserFactory) CreateParser(input [][]byte, grammar int) scan.Parser {
	return &scan.ParserData{
		Input:         input,
		ParsingResult: make([]scan.Peer, 0),
		Grammar:       grammar,
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *ParserFactory) IsInterfaceNil() bool {
	return pf == nil
}
