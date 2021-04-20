package trieIterators

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/node/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/stretchr/testify/require"
)

func TestValidatorUnstakedValue(t *testing.T) {
	args := createMockArgs()

	args.QueryService = &mock.SCQueryServiceStub{
		ExecuteQueryCalled: func(query *process.SCQuery) (*vmcommon.VMOutput, error) {
			return &vmcommon.VMOutput{
				ReturnCode: vmcommon.Ok,
				ReturnData: [][]byte{{1}, {2}, {1}, {2}},
			}, nil
		},
	}
	commonProcessor := &commonStakingProcessor{
		queryService: args.QueryService,
	}

	unstakedValue, _ := commonProcessor.getValidatorUnstakedValue([]byte("validator"))
	require.Equal(t, big.NewInt(2), unstakedValue)
}
