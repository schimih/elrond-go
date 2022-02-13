package scan

import (
	"testing"

	core "github.com/elrond-go/cmd/vat/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNmapScanner_ScanFail(t *testing.T) {
	scanner := &NmapScanner{
		Name:   "testName",
		Target: "testTarget",
		Status: core.NOT_STARTED,
		Cmd:    "testCmd",
	}

	res, _ := scanner.Scan()

	assert.Nil(t, res)
	require.Equal(t, core.IN_PROGRESS, scanner.Status)
}

func TestNmapScanner_ScanPass(t *testing.T) {
	scanner := &NmapScanner{
		Name:   "testName",
		Target: "128.199.37.240",
		Status: core.NOT_STARTED,
		Cmd:    "nmap -Pn -p22",
	}

	_, err := scanner.Scan()

	assert.Nil(t, err)
	require.Equal(t, core.FINISHED, scanner.Status)
}
