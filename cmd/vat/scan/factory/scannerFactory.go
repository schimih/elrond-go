package factory

import (
	"fmt"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type ScannerFactory struct {
}

func NewScannerFactory() *ScannerFactory {
	return &ScannerFactory{}
}

func (factory *ScannerFactory) CreateScanner(target string, analysisType utils.AnalysisType) scan.Scanner {
	switch analysisType {
	case utils.TCP_ELROND:
		return &scan.NmapScanner{
			Name:   "TCP-ELROND",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, scan.NMAP_TCP_ELROND),
		}
	case utils.TCP_WEB:
		return &scan.NmapScanner{
			Name:   "TCP-WEB",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, scan.NMAP_TCP_WEB),
		}
	case utils.TCP_SSH_ALGOS:
		return &scan.NmapScanner{
			Name:   "TCP-SSH-ALGO",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    fmt.Sprintf("nmap --script ssh2-enum-algos %s -oX -", target),
		}
	case utils.TCP_STANDARD:
		return &scan.NmapScanner{
			Name:   "TCP-TOP",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, scan.NMAP_TCP_STANDARD),
		}
	case utils.TCP_REQ1:
		return &scan.NmapScanner{
			Name:   "TCP-REQ",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, scan.NMAP_TCP_REQ1),
		}
	case utils.TCP_POLITE_REQ1:
		return &scan.PoliteScanner{
			Host: target,
		}
	case utils.TCP_BRUTE_REQ1:
		return &scan.NmapScanner{
			Name:   "TCP-BRUTE-SSH",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, scan.NMAP_TCP_REQ1),
		}
	default:
		return nil
	}
}

func constructCmd(target string, args scan.NmapCommand) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *ScannerFactory) IsInterfaceNil() bool {
	return factory == nil
}
