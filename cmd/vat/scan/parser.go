package scan

import (
	"strings"

	core "github.com/elrond-go/cmd/vat/core"
	gonmap "github.com/lair-framework/go-nmap"
)

type ParserData struct {
	Input             [][]byte
	AnalyzedTargets   []ScannedTarget
	Grammar           core.AnalysisType
	SlicedParsedInput []*gonmap.NmapRun
}

func (pD *ParserData) Parse() (parsingResults []ScannedTarget) {
	for _, byteArray := range pD.Input {
		parsedNmapResult, err := gonmap.Parse(byteArray)
		if err != nil {
			if !strings.Contains(err.Error(), "exit status 1") {
				log.Error(err.Error())
			}
		}
		pD.SlicedParsedInput = append(pD.SlicedParsedInput, parsedNmapResult)
	}

	pD.translateInput()

	return pD.AnalyzedTargets
}

func (pD *ParserData) translateInput() {
	for _, nmapRun := range pD.SlicedParsedInput {
		for hidx, host := range nmapRun.Hosts {
			pD.translateTarget(hidx, host)
		}
	}
}

func (pD *ParserData) translateTarget(id int, host gonmap.Host) {
	pS := createPortSlice(host)
	translatedPortSlice := pS.translatePortSlice()
	analyzedTarget := NewScannedTarget(uint(id), host.Addresses[0].Addr, translatedPortSlice, core.SCANNED, pD.Grammar)
	pD.AnalyzedTargets = append(pD.AnalyzedTargets, analyzedTarget)
}

func CreateParser(input [][]byte, grammar core.AnalysisType) Parser {
	return &ParserData{
		Input:             input,
		AnalyzedTargets:   make([]ScannedTarget, 0),
		Grammar:           grammar,
		SlicedParsedInput: make([]*gonmap.NmapRun, 0),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pD *ParserData) IsInterfaceNil() bool {
	return pD == nil
}
