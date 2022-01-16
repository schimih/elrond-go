package factory

import (
	"fmt"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/vat"
)

type NmapScannerFactory struct {
}

func NewNmapScannerFactory() *NmapScannerFactory {
	return &NmapScannerFactory{}
}

func (factory *NmapScannerFactory) CreateScanner(target string, analysisType vat.AnalysisType) scan.Scanner {
	switch analysisType {
	case vat.TCP_ELROND:
		return &scan.ArgNmapScanner{
			Name:   "TCP-ELROND",
			Target: target,
			Status: vat.NOT_STARTED,
			Cmd:    constructCmd(target, vat.NMAP_TCP_ELROND),
		}
	case vat.TCP_WEB:
		return &scan.ArgNmapScanner{
			Name:   "TCP-WEB",
			Target: target,
			Status: vat.NOT_STARTED,
			Cmd:    constructCmd(target, vat.NMAP_TCP_WEB),
		}
	case vat.TCP_SSH:
		return &scan.ArgNmapScanner{
			Name:   "TCP-SSH",
			Target: target,
			Status: vat.NOT_STARTED,
			Cmd:    constructCmd(target, vat.NMAP_TCP_SSH),
		}
	case vat.TCP_STANDARD:
		return &scan.ArgNmapScanner{
			Name:   "TCP-TOP",
			Target: target,
			Status: vat.NOT_STARTED,
			Cmd:    constructCmd(target, vat.NMAP_TCP_STANDARD),
		}
	default:
		return nil
	}
}

func constructCmd(target string, args vat.NmapCommand) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *NmapScannerFactory) IsInterfaceNil() bool {
	return factory == nil
}
