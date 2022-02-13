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
	"github.com/elrond-go/cmd/vat/export"
	"github.com/elrond-go/cmd/vat/scan/factory"
)

var log = logger.GetOrCreate("vat")

func main() {
	log.Info("Starting VAT")
	var err error

	// Start seednode
	messenger, err := startSeedNode()
	if err != nil {
		log.Error("could not start seed node", "error", err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Analysis is now running")
	mainLoop(messenger, sigs)

}

func mainLoop(messenger p2p.Messenger, stop chan os.Signal) {
	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()

	sF := factory.NewScannerFactory()
	pF := factory.NewParserFactory()
	d := analysis.NewP2pDiscoverer(messenger)
	a, _ := analysis.NewAnalyzer(d, sF, pF)

	eF := evaluation.NewEvaluatorFactory()
	report := evaluation.NewReport(eF, sF)

	fF := export.NewFormatterFactory()
	vatExport, _ := export.NewAnalysisManager(fF)

	for {
		interval := time.Second * 5

		analyzedTargets := a.StartJob(vatExport.AnalysisType)
		evaluatedTargets, _ := report.RunEvaluation(analyzedTargets, vatExport.EvaluationType)
		vatExport.CompleteRound(evaluatedTargets)

		timer.Reset(interval)

		select {
		case <-timer.C:
		case <-stop:
			log.Info("terminating at user's signal...")
			return
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

func createNetworkMessenger(p2pConfig config.P2PConfig, marshalizer marshal.Marshalizer) (p2p.Messenger, error) {
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

	generalConfig, err := loadMainConfig("./cmd/vat/config/config.toml")
	if err != nil {
		return nil, err
	}

	internalMarshalizer, err := factoryMarshalizer.NewMarshalizer(generalConfig.Marshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("error creating marshalizer (internal): %s", err.Error())
	}

	log.Info("Starting Seed Node")

	p2pConfig, err := common.LoadP2PConfig("./cmd/vat/config/p2p.toml")
	if err != nil {
		return nil, err
	}

	messenger, err = createNetworkMessenger(*p2pConfig, internalMarshalizer)
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
