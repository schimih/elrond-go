package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	factoryMarshalizer "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/epochStart/bootstrap/disabled"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/elrond-go/cmd/vat/analysis"
	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/scan/factory"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/urfave/cli"
)

type cfg struct {
	vatAnalysisType string
}

const (
	defaultLogsPath     = "logs"
	logFilePrefix       = "elrond-seed"
	filePathPlaceholder = "[path]"
)

var (
	vatTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}
 VERSION:
	{{.Version}}
	{{end}}
 `

	analysisType = cli.StringFlag{
		Name:        "analysis-type",
		Usage:       "Provide type of analysis. Default full test",
		Value:       "full",
		Destination: &argsConfig.vatAnalysisType,
	}
	argsConfig           = &cfg{}
	p2pConfigurationFile = "./config/p2p.toml"
)

const (
	TCP_ELROND = iota
	TCP_OUTSIDE_ELROND
	TCP_WEB
	TCP_SSH
	TCP_FULL
	TCP_STANDARD
)

var log = logger.GetOrCreate("vat")

func main() {
	log.Info("Starting VAT")
	app := cli.NewApp()
	cli.AppHelpTemplate = vatTemplate
	app.Name = "Vulnerability Analysis Tool"
	app.Version = "v0.0.1"
	app.Usage = "This tool will be used for security checks on Elrond EcoSystem (v0.0.1 - portscanner and ssh access)"
	app.Flags = []cli.Flag{
		analysisType,
	}

	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		return startVulnerabilityAnalysis(ctx)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startVulnerabilityAnalysis(ctx *cli.Context) error {
	var err error
	// to implement ctx for test flag

	// Start seednode
	messenger, err := startSeedNode()
	if err != nil {
		return err
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Analysis is now running")
	analysisType := utils.TCP_WEB
	mainLoop(messenger, sigs, analysisType)

	return nil
}

func mainLoop(messenger p2p.Messenger, stop chan os.Signal, analysisType utils.AnalysisType) {
	sf := factory.NewNmapScannerFactory()
	pf := factory.NewParserFactory()
	ef := evaluation.NewEvaluatorFactory()
	d := analysis.NewP2pDiscoverer(messenger)
	a, _ := analysis.NewAnalyzer(d, sf, pf, analysisType)

	report := evaluation.NewEvaluationReport(ef)

	//manager := manager.NewAnalysisManager(analysisType)
	evaluationType := utils.PortStatusEvaluation
	for {
		select {
		case <-stop:
			log.Info("terminating at user's signal...")
			return
		case <-time.After(time.Second * 5):
			a.DiscoverTargets()
			analyzedTargets := a.AnalyzeNewlyDiscoveredTargets()
			report.StartEvaluation(analyzedTargets, evaluationType)
			report.DisplayToTable()
			log.Info("Added targets", "targets", len(a.DiscoveredTargets))
		}
	}
}

func loadMainConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createNode(p2pConfig config.P2PConfig, marshalizer marshal.Marshalizer) (p2p.Messenger, error) {
	arg := libp2p.ArgsNetworkMessenger{
		Marshalizer:          marshalizer,
		ListenAddress:        libp2p.ListenAddrWithIp4AndTcp,
		P2pConfig:            p2pConfig,
		SyncTimer:            &libp2p.LocalSyncTimer{},
		PreferredPeersHolder: disabled.NewPreferredPeersHolder(),
		NodeOperationMode:    p2p.NormalOperation,
	}

	return libp2p.NewNetworkMessenger(arg)
}

func startSeedNode() (messenger p2p.Messenger, err error) {

	generalConfig, err := loadMainConfig("./config/config.toml")
	if err != nil {
		return nil, err
	}

	internalMarshalizer, err := factoryMarshalizer.NewMarshalizer(generalConfig.Marshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("error creating marshalizer (internal): %s", err.Error())
	}

	log.Info("Starting Seed Node")

	p2pConfig, err := common.LoadP2PConfig(p2pConfigurationFile)
	if err != nil {
		return nil, err
	}

	messenger, err = createNode(*p2pConfig, internalMarshalizer)
	if err != nil {
		return nil, err
	}

	log.Info("Starting Bootstrap")
	err = messenger.Bootstrap()
	if err != nil {
		return nil, err
	}

	return messenger, nil
}
