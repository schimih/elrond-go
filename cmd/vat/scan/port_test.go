package scan

import (
	"testing"

	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestNewPort(t *testing.T) {
	expectedPort := Port{
		ID:         0,
		Number:     10000,
		Protocol:   "testProtocol",
		State:      utils.Closed,
		Owner:      "testOwner",
		Type:       utils.ElrondPort,
		Importance: utils.JudgementNoRisk,
	}

	newPort := NewPort(0, 10000, "testProtocol", utils.Closed, "testOwner", utils.ElrondPort)

	require.Equal(t, expectedPort, newPort)
}
