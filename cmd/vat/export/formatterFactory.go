package export

import (
	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/elrond-go/cmd/vat/utils"
)

type FormatterFactory struct {
}

func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

func (fF *FormatterFactory) CreateFormatter(formatType utils.OutputType) (formatter Formatter, err error) {
	switch formatType {
	case utils.Table:
		return &TableFormatter{
			header:    make([]string, 0),
			dataLines: make([]*display.LineData, 0),
		}, nil
	case utils.JSON:
		return &JsonFormatter{}, nil
	case utils.XML:
		return &XMLFormatter{}, nil
	default:
		return nil, utils.ErrNoFormatterType
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (fF *FormatterFactory) IsInterfaceNil() bool {
	return fF == nil
}
