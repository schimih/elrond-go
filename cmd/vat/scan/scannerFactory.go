package scan

import (
	"fmt"

	core "github.com/elrond-go/cmd/vat/core"
)

type ScannerFactory struct {
}

func NewScannerFactory() *ScannerFactory {
	return &ScannerFactory{}
}

func (factory *ScannerFactory) CreateScanner(target string, analysisType core.AnalysisType) Scanner {
	switch analysisType {
	case core.TCP_ELROND:
		return &NmapScanner{
			Name:   "TCP-ELROND",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_ELROND),
		}
	case core.TCP_WEB:
		return &NmapScanner{
			Name:   "TCP-WEB",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_WEB),
		}
	case core.TCP_SSH_ALGOS:
		return &NmapScanner{
			Name:   "TCP-SSH-ALGO",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    fmt.Sprintf("nmap --script ssh2-enum-algos %s -oX -", target),
		}
	case core.TCP_STANDARD:
		return &NmapScanner{
			Name:   "TCP-TOP",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_STANDARD),
		}
	case core.TCP_REQ1:
		return &NmapScanner{
			Name:   "TCP-REQ",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_REQ1),
		}
	case core.TCP_POLITE_REQ1:
		return &PoliteScanner{
			Host: target,
		}
	case core.TCP_BRUTE_REQ1:
		return &NmapScanner{
			Name:   "TCP-BRUTE-SSH",
			Target: target,
			Status: core.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_REQ1),
		}
	default:
		return nil
	}
}

func constructCmd(target string, args NmapCommand) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *ScannerFactory) IsInterfaceNil() bool {
	return factory == nil
}
