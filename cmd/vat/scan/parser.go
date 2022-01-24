package scan

import (
	"strings"

	"github.com/elrond-go/cmd/vat/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type ParserData struct {
	Input             [][]byte
	AnalyzedTargets   []ScannedTarget
	Grammar           utils.AnalysisType
	SlicedParsedInput []*go_nmap.NmapRun
}

func (pD *ParserData) Parse() (parsingResults []ScannedTarget) {
	for _, byteArray := range pD.Input {
		parsedNmapResult, err := go_nmap.Parse(byteArray)
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

func (pD *ParserData) translateTarget(id int, host go_nmap.Host) {
	pS := createPortSlice(host)
	translatedPortSlice := pS.translatePortSlice()
	analyzedTarget := NewScannedTarget(uint(id), host.Addresses[0].Addr, translatedPortSlice, host.Status.State, utils.SCANNED, pD.Grammar)
	pD.AnalyzedTargets = append(pD.AnalyzedTargets, analyzedTarget)
}

// IsInterfaceNil returns true if there is no value under the interface
func (pD *ParserData) IsInterfaceNil() bool {
	return pD == nil
}
