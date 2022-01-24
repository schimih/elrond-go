package manager

import (
	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/elrond-go/cmd/vat/utils"
)

type FormatterFactory struct {
}

func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

func (fF *FormatterFactory) CreateFormatter(formatType utils.OutputType) Formatter {
	switch formatType {
	case utils.Table:
		return &TableFormatter{
			header:    make([]string, 0),
			dataLines: make([]*display.LineData, 0),
		}
	case utils.JSON:
		return &JsonFormatter{}
	case utils.XML:
		return &XMLFormatter{}
	default:
		return nil
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (fF *FormatterFactory) IsInterfaceNil() bool {
	return fF == nil
}
