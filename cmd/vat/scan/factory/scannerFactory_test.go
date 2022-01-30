package factory

import (
	"testing"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateScanner(t *testing.T) {
	sF := NewScannerFactory()

	scanner := sF.CreateScanner("testTarget", utils.TCP_ELROND)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_WEB)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_SSH_ALGOS)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_STANDARD)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_POLITE_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.PoliteScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", utils.TCP_BRUTE_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &scan.NmapScanner{}, scanner)
}
