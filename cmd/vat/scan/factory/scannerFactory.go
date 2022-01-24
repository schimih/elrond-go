package factory

import (
	"fmt"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type NmapScannerFactory struct {
}

func NewNmapScannerFactory() *NmapScannerFactory {
	return &NmapScannerFactory{}
}

func (factory *NmapScannerFactory) CreateScanner(target string, analysisType utils.AnalysisType) scan.Scanner {
	switch analysisType {
	case utils.TCP_ELROND:
		return &scan.ArgNmapScanner{
			Name:   "TCP-ELROND",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, utils.NMAP_TCP_ELROND),
		}
	case utils.TCP_WEB:
		return &scan.ArgNmapScanner{
			Name:   "TCP-WEB",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, utils.NMAP_TCP_WEB),
		}
	case utils.TCP_SSH:
		return &scan.ArgNmapScanner{
			Name:   "TCP-SSH",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, utils.NMAP_TCP_SSH),
		}
	case utils.TCP_STANDARD:
		return &scan.ArgNmapScanner{
			Name:   "TCP-TOP",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, utils.NMAP_TCP_STANDARD),
		}
	case utils.TCP_REQ1:
		return &scan.ArgNmapScanner{
			Name:   "TCP-REQ",
			Target: target,
			Status: utils.NOT_STARTED,
			Cmd:    constructCmd(target, utils.NMAP_TCP_REQ1),
		}
	default:
		return nil
	}
}

func constructCmd(target string, args utils.NmapCommand) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *NmapScannerFactory) IsInterfaceNil() bool {
	return factory == nil
}
