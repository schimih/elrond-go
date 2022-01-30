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
	return &scan.NmapScanner{
		Name:   "TCP-SSH",
		Target: target,
		Status: utils.NOT_STARTED,
		Cmd:    "Test_Cmd_Should_Fail"}
}

func (fpf *FakeParserFactory) CreateParser(input [][]byte, grammar utils.AnalysisType) scan.Parser {
	return &scan.ParserData{
		Input:           input,
		AnalyzedTargets: make([]scan.ScannedTarget, 0),
		Grammar:         grammar,
	}
}

func TestNewAnalyzer(t *testing.T) {
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, sff, fpf)
	assert.False(t, check.IfNil(na))
	assert.Nil(t, err)
}

func TestNewAnalyzer_DiscovererNilCheck(t *testing.T) {
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(nil, sff, fpf)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "Discoverer needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ScannerFactoryNilCheck(t *testing.T) {
	fd := &FakeDiscoverer{}
	//sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, nil, fpf)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ScannerFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestNewAnalyzer_ParserFactoryNilCheck(t *testing.T) {
	fd := &FakeDiscoverer{}
	sff := &FakeScannerFactory{}
	//fpf := &FakeParserFactory{}
	na, err := NewAnalyzer(fd, sff, nil)
	assert.True(t, check.IfNil(na))
	expectedErrorString := "ParserFactory needed"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestAnalyzer_DiscoverNewPeers(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf)
	discovererStub.DiscoverNewTargetsCalled = func(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget) {
		return make([]DiscoveredTarget, 2)
	}
	na.DiscoverTargets()

	require.Equal(t, 2, len(na.DiscoveredTargets))
}

func TestAnalyzeNewlyDiscoveredTargets(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf)
	analysisType := utils.TCP_WEB
	na.AnalyzeNewlyDiscoveredTargets(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsNew(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf)
	analysisType := utils.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         utils.NEW,
	}
	na.DiscoveredTargets = append(na.DiscoveredTargets, DiscoveredTarget)
	na.AnalyzeNewlyDiscoveredTargets(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsExpired(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf)
	analysisType := utils.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         utils.EXPIRED,
	}
	na.DiscoveredTargets = append(na.DiscoveredTargets, DiscoveredTarget)
	na.AnalyzeNewlyDiscoveredTargets(analysisType)
}

func TestAnalyzeNewlyDiscoveredTargets_ActualStatusIsNorNewOrExpired(t *testing.T) {
	discovererStub := NewDiscovererStub()
	sff := &FakeScannerFactory{}
	fpf := &FakeParserFactory{}
	na, _ := NewAnalyzer(discovererStub, sff, fpf)
	analysisType := utils.TCP_WEB
	DiscoveredTarget := DiscoveredTarget{
		ID:             0,
		Protocol:       "Test_Protocol",
		Address:        "Test_Address",
		ConnectionPort: "Test_Port",
		Status:         utils.SCANNED,
	}
	na.DiscoveredTargets = append(na.DiscoveredTargets, DiscoveredTarget)
	na.AnalyzeNewlyDiscoveredTargets(analysisType)
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
