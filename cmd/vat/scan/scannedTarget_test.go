package scan

import (
	"testing"

	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestNewScannedTarget(t *testing.T) {
	expected := ScannedTarget{
		ID:           0,
		Protocol:     "",
		Address:      "testAddress",
		Standard:     "",
		Ports:        make([]Port, 0),
		Architecture: "",
		Status:       utils.NEW,
		AnalysisType: utils.TCP_SSH,
	}

	newScannedTarget := NewScannedTarget(0, "testAddress", make([]Port, 0), utils.NEW, utils.TCP_SSH)

	require.Equal(t, expected, newScannedTarget)
}
