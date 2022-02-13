package analysis

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/elrond-go/cmd/vat/core"
	"github.com/elrond-go/cmd/vat/scan"
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

func (sff *FakeScannerFactory) CreateScanner(target string, analysisType core.AnalysisType) (Scanner scan.Scanner) {
	return &scan.NmapScanner{
		Name:   "TCP-SSH",
		Target: target,
		Status: core.NOT_STARTED,
		Cmd:    "Test_Cmd_Should_Fail"}
}

func (fpf *FakeParserFactory) CreateParser(input [][]byte, grammar core.AnalysisType) scan.Parser {
	return &scan.ParserData{
		Input:           input,
		AnalyzedTargets: make([]scan.ScannedTarget, 0),
		Grammar:         grammar,
	}
}

func TestNewAnalyzer(t *testing.T) {
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	na, err := NewAnalyzer(fd, sff)
	assert.False(t, check.IfNil(na))
	assert.Nil(t, err)
}

func TestNewAnalyzer_DiscovererNilCheck(t *testing.T) {
	sff := &FakeScannerFactory{}
	na, err := NewAnalyzer(nil, sff)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "Discoverer needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ScannerFactoryNilCheck(t *testing.T) {
	fd := &FakeDiscoverer{}
	na, err := NewAnalyzer(fd, nil)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ScannerFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ParserFactoryNilCheck(t *testing.T) {
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	na, err := NewAnalyzer(fd, sff)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ParserFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestAnalyzer_DiscoverNewPeers(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	na, _ := NewAnalyzer(discovererStub, sff)
	discovererStub.DiscoverNewTargetsCalled = func(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget) {
		return make([]DiscoveredTarget, 2)
	}
	na.discoverTargets()

	require.Equal(t, 2, len(na.discoveredTargets))
}

func TestAnalyzeNewlyDiscoveredTargets(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	na, _ := NewAnalyzer(discovererStub, sff)
	analysisType := core.TCP_WEB
	na.StartJob(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsNew(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	na, _ := NewAnalyzer(discovererStub, sff)
	analysisType := core.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         core.NEW,
	}
	na.discoveredTargets = append(na.discoveredTargets, DiscoveredTarget)
	na.StartJob(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsExpired(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	na, _ := NewAnalyzer(discovererStub, sff)
	analysisType := core.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         core.EXPIRED,
	}
	na.discoveredTargets = append(na.discoveredTargets, DiscoveredTarget)
	na.StartJob(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsNorNewOrExpired(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	na, _ := NewAnalyzer(discovererStub, sff)
	analysisType := core.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         core.SCANNED,
	}
	na.discoveredTargets = append(na.discoveredTargets, DiscoveredTarget)
	na.StartJob(analysisType)
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
