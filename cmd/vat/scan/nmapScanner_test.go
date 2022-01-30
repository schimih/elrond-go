package scan

import (
	"testing"

	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNmapScanner_ScanFail(t *testing.T) {
	scanner := &NmapScanner{
		Name:   "testName",
		Target: "testTarget",
		Status: utils.NOT_STARTED,
		Cmd:    "testCmd",
	}

	res, _ := scanner.Scan()

	assert.Nil(t, res)
	require.Equal(t, utils.IN_PROGRESS, scanner.Status)
}

func TestNmapScanner_ScanPass(t *testing.T) {
	scanner := &NmapScanner{
		Name:   "testName",
		Target: "128.199.37.240",
		Status: utils.NOT_STARTED,
		Cmd:    "nmap -Pn -p22",
	}

	_, err := scanner.Scan()

	assert.Nil(t, err)
	require.Equal(t, utils.FINISHED, scanner.Status)
}
