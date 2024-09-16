package scan

import (
	"testing"

	core "github.com/elrond-go/cmd/vat/core"
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
		Status:       core.NEW,
		AnalysisType: core.TCP_SSH,
	}

	newScannedTarget := NewScannedTarget(0, "testAddress", make([]Port, 0), core.NEW, core.TCP_SSH)

	require.Equal(t, expected, newScannedTarget)
}
