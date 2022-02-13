package scan

import (
	"testing"

	core "github.com/elrond-go/cmd/vat/core"
	"github.com/stretchr/testify/require"
)

func TestNewPort(t *testing.T) {
	expectedPort := Port{
		ID:         0,
		Number:     10000,
		Protocol:   "testProtocol",
		State:      core.Closed,
		Owner:      "testOwner",
		Type:       core.ElrondPort,
		Importance: core.JudgementNoRisk,
	}

	newPort := NewPort(0, 10000, "testProtocol", core.Closed, "testOwner", core.ElrondPort)

	require.Equal(t, expectedPort, newPort)
}
