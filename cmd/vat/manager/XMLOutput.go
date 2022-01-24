package manager

type XMLFormatter struct {
}

func (xF *XMLFormatter) Output(rankedReport RankedReport) {

}

// IsInterfaceNil returns true if there is no value under the interface
func (xF *XMLFormatter) IsInterfaceNil() bool {
	return xF == nil
}
