package scan

import (
	"strings"

	"github.com/elrond-go/cmd/vat/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type ParserData struct {
	Input           [][]byte
	AnalyzedTargets []AnalyzedTarget
	Grammar         utils.AnalysisType
}

func (p *ParserData) Parse() (parsingResults []AnalyzedTarget) {
	slicedParserInput := make([]*go_nmap.NmapRun, 0)
	for _, byteArray := range p.Input {
		parsedNmapResult, err := go_nmap.Parse(byteArray)
		if err != nil {
			if !strings.Contains(err.Error(), "exit status 1") {
				log.Error(err.Error())
			}
		}
		slicedParserInput = append(slicedParserInput, parsedNmapResult)
	}
	p.translateInput(slicedParserInput)
	return p.AnalyzedTargets
}

func (p *ParserData) translateInput(NmapScanResult []*go_nmap.NmapRun) {
	for _, nmapRun := range NmapScanResult {
		for hidx, host := range nmapRun.Hosts {
			p.translateTarget(hidx, host)
		}
	}
}

func (p *ParserData) translateTarget(id int, host go_nmap.Host) {
	pS := createPortSlice(host)
	analyzedTarget := NewAnalyzedTarget(uint(id), host.Addresses[0].Addr, pS.translatePortSlice(), host.Status.State, utils.SCANNED, p.Grammar)
	p.AnalyzedTargets = append(p.AnalyzedTargets, analyzedTarget)
}

// IsInterfaceNil returns true if there is no value under the interface
func (p *ParserData) IsInterfaceNil() bool {
	return p == nil
}
