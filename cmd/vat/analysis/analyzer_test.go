package analysis

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeDiscoverer struct {
}

type FakeParserFactory struct {
}

type FakeScannerFactory struct {
}

func (fd *FakeDiscoverer) DiscoverNewTargets(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget) {
	targets = existingTargets

	return
}

func (sff *FakeScannerFactory) CreateScanner(target string, analysisType utils.AnalysisType) (Scanner scan.Scanner) {
	return &scan.ArgNmapScanner{Name: "TCP-SSH",
		Target: target,
		Status: utils.NOT_STARTED,
		Cmd:    "constructCmd(target, NMAP_TCP_SSH)"}
}

func (fpf *FakeParserFactory) CreateParser(input [][]byte, grammar utils.AnalysisType) scan.Parser {
	return &scan.ParserData{
		Input:           input,
		AnalyzedTargets: make([]scan.AnalyzedTarget, 0),
		Grammar:         grammar,
	}
}

func TestNewAnalyzer(t *testing.T) {
	analysisType := utils.TCP_ELROND
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, sff, fpf, analysisType)
	assert.False(t, check.IfNil(na))
	assert.Nil(t, err)
}

func TestNewAnalyzer_DiscovererNilCheck(t *testing.T) {
	analysisType := utils.TCP_ELROND
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(nil, sff, fpf, analysisType)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "Discoverer needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ScannerFactoryNilCheck(t *testing.T) {
	analysisType := utils.TCP_ELROND
	fd := &FakeDiscoverer{}
	//sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, nil, fpf, analysisType)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ScannerFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ParserFactoryNilCheck(t *testing.T) {
	analysisType := utils.TCP_ELROND
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	//fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, sff, nil, analysisType)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ParserFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestAnalyzer_DiscoverNewPeers(t *testing.T) {
	analysisType := utils.TCP_ELROND
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf, analysisType)
	discovererStub.DiscoverNewTargetsCalled = func(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget) {
		return make([]DiscoveredTarget, 2)
	}
	na.DiscoverTargets()

	require.Equal(t, 2, len(na.DiscoveredTargets))
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
