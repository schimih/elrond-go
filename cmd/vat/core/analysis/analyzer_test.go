package analysis

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/elrond-go/cmd/vat/core/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeDiscoverer struct {
}

type FakeParserFactory struct {
}

type FakeScannerFactory struct {
}

func (fd *FakeDiscoverer) DiscoverNewTargets(existingTargets []Target) (targets []Target) {
	targets = existingTargets

	return
}

func (sff *FakeScannerFactory) CreateScanner(target string, analysisType int) (Scanner scan.Scanner) {
	return &scan.ArgNmapScanner{Name: "TCP-SSH",
		Target: target,
		Status: scan.NOT_STARTED,
		Cmd:    "constructCmd(target, NMAP_TCP_SSH)"}
}

func (fpf *FakeParserFactory) CreateParser(input [][]byte, grammar int) scan.Parser {
	return &scan.ParserData{
		Input:         input,
		ParsingResult: make([]scan.Peer, 0),
		Grammar:       grammar,
	}
}

func TestNewAnalyzer(t *testing.T) {
	analysisType := 0
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, sff, fpf, analysisType)
	assert.False(t, check.IfNil(na))
	assert.Nil(t, err)
}

func TestAnalyzer_DiscoverNewPeers(t *testing.T) {
	analysisType := 0
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf, analysisType)
	discovererStub.DiscoverNewTargetsCalled = func(existingTargets []Target) (targets []Target) {
		return make([]Target, 2)
	}
	na.DiscoverNewPeers()

	require.Equal(t, 2, len(na.Targets))
}

// IsInterfaceNil returns true if there is no value under the interface
func (fpf *FakeParserFactory) IsInterfaceNil() bool {
	return fpf == nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (fsf *FakeScannerFactory) IsInterfaceNil() bool {
	return fsf == nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *FakeDiscoverer) IsInterfaceNil() bool {
	return d == nil
}
