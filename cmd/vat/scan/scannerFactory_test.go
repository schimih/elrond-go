package scan

import (
	"testing"

	core "github.com/elrond-go/cmd/vat/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateScanner(t *testing.T) {
	sF := NewScannerFactory()

	scanner := sF.CreateScanner("testTarget", core.TCP_ELROND)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_WEB)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_SSH_ALGOS)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_STANDARD)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_POLITE_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &PoliteScanner{}, scanner)

	scanner = sF.CreateScanner("testTarget", core.TCP_BRUTE_REQ1)
	assert.NotNil(t, scanner)
	assert.IsType(t, &NmapScanner{}, scanner)
}
