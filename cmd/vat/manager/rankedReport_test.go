package manager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRankedReport(t *testing.T) {
	rankedReport := NewRankedReport()

	require.False(t, rankedReport.IsInterfaceNil())
}
