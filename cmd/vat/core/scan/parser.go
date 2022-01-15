package scan

import (
	"strings"

	go_nmap "github.com/lair-framework/go-nmap"
)

type ParserData struct {
	Input         [][]byte
	ParsingResult []Peer
	Grammar       int
}

type TargetStatus string

const (
	NEW     TargetStatus = "NEW"
	SCANNED TargetStatus = "SCANNED"
	EXPIRED TargetStatus = "EXPIRED"
)

func (p *ParserData) Parse() (parsingResults []Peer) {
	nmapResultSlice := make([]*go_nmap.NmapRun, 0)
	for _, byteArray := range p.Input {
		nmapResult, err := go_nmap.Parse(byteArray)
		if err != nil {
			if !strings.Contains(err.Error(), "exit status 1") {
				log.Error(err.Error())
			}
		}
		nmapResultSlice = append(nmapResultSlice, nmapResult)
	}
	p.processInput(nmapResultSlice)
	return p.ParsingResult
}

func (p *ParserData) processInput(NmapScanResult []*go_nmap.NmapRun) {
	for _, nmapRun := range NmapScanResult {
		for hidx, host := range nmapRun.Hosts {
			p.process(hidx, host)
		}
	}
}

func (p *ParserData) process(id int, host go_nmap.Host) {
	pS := createPortSlice(host)
	peer := NewPeer(uint(id), host.Addresses[0].Addr, pS.translatePortSlice(), host.Status.State, SCANNED, p.Grammar)
	p.ParsingResult = append(p.ParsingResult, peer)
}

// IsInterfaceNil returns true if there is no value under the interface
func (p *ParserData) IsInterfaceNil() bool {
	return p == nil
}
