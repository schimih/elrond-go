package detector_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data/block"
	mock2 "github.com/ElrondNetwork/elrond-go/epochStart/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/block/interceptedBlocks"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/process/slash"
	"github.com/ElrondNetwork/elrond-go/process/slash/detector"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewHeaderSlashingDetector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() sharding.NodesCoordinator
		expectedErr error
	}{
		{
			args: func() sharding.NodesCoordinator {
				return nil
			},
			expectedErr: process.ErrNilShardCoordinator,
		},
		{
			args: func() sharding.NodesCoordinator {
				return &mock.NodesCoordinatorMock{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := detector.NewHeaderSlashingDetector(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestHeaderSlashingDetector_VerifyData_CannotCastData_ExpectError(t *testing.T) {
	t.Parallel()

	sd, _ := detector.NewHeaderSlashingDetector(&mock.NodesCoordinatorMock{})

	res, err := sd.VerifyData(&testscommon.InterceptedDataStub{})

	require.Nil(t, res)
	require.Equal(t, process.ErrCannotCastInterceptedDataToHeader, err)
}

func TestHeaderSlashingDetector_VerifyData_CannotGetProposer_ExpectError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("cannot get proposer")
	sd, _ := detector.NewHeaderSlashingDetector(&mock2.NodesCoordinatorStub{
		ComputeConsensusGroupCalled: func(_ []byte, _ uint64, _ uint32, _ uint32) ([]sharding.Validator, error) {
			return nil, expectedErr
		},
	})

	res, err := sd.VerifyData(&interceptedBlocks.InterceptedHeader{})

	require.Nil(t, res)
	require.Equal(t, expectedErr, err)
}

func TestHeaderSlashingDetector_VerifyData_NoSlashing(t *testing.T) {
	t.Parallel()

	sd, _ := detector.NewHeaderSlashingDetector(&mock2.NodesCoordinatorStub{
		ComputeConsensusGroupCalled: func(_ []byte, _ uint64, _ uint32, _ uint32) ([]sharding.Validator, error) {
			return []sharding.Validator{mock.NewValidatorMock([]byte("proposer1"))}, nil
		},
	})

	hData := createInterceptedHeaderData(2, []byte("seed"))
	res, _ := sd.VerifyData(hData)
	require.Equal(t, res.GetType(), slash.None)
	require.Equal(t, res.GetLevel(), slash.Level0)

	res, _ = sd.VerifyData(hData)
	require.Equal(t, res.GetType(), slash.None)
	require.Equal(t, res.GetLevel(), slash.Level0)

	res, _ = sd.VerifyData(hData)
	require.Equal(t, res.GetType(), slash.None)
	require.Equal(t, res.GetLevel(), slash.Level0)
}

func TestHeaderSlashingDetector_VerifyData_MultipleProposal(t *testing.T) {
	t.Parallel()

	sd, _ := detector.NewHeaderSlashingDetector(&mock2.NodesCoordinatorStub{
		ComputeConsensusGroupCalled: func(_ []byte, _ uint64, _ uint32, _ uint32) ([]sharding.Validator, error) {
			return []sharding.Validator{mock.NewValidatorMock([]byte("proposer1"))}, nil
		},
	})

	hData1 := createInterceptedHeaderData(2, []byte("seed1"))
	tmp, _ := sd.VerifyData(hData1)

	require.Equal(t, tmp.GetType(), slash.None)
	require.Equal(t, tmp.GetLevel(), slash.Level0)

	hData2 := createInterceptedHeaderData(2, []byte("seed2"))
	tmp, _ = sd.VerifyData(hData2)
	res := tmp.(slash.MultipleProposalProofHandler)

	require.Equal(t, res.GetType(), slash.MultipleProposal)
	require.Equal(t, res.GetLevel(), slash.Level1)
	require.Len(t, res.GetHeaders(), 2)
	require.Equal(t, res.GetHeaders()[0], hData1)
	require.Equal(t, res.GetHeaders()[1], hData2)

	hData3 := createInterceptedHeaderData(2, []byte("seed3"))
	tmp, _ = sd.VerifyData(hData3)
	res = tmp.(slash.MultipleProposalProofHandler)

	require.Equal(t, res.GetType(), slash.MultipleProposal)
	require.Equal(t, res.GetLevel(), slash.Level2)
	require.Len(t, res.GetHeaders(), 3)
	require.Equal(t, res.GetHeaders()[0], hData1)
	require.Equal(t, res.GetHeaders()[1], hData2)
	require.Equal(t, res.GetHeaders()[2], hData3)
}

func TestHeaderSlashingDetector_ValidateProof_SimpleSlashingProof_DifferentSlashLevelsAndTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (slash.SlashingLevel, slash.SlashingType)
		expectedErr error
	}{
		{
			args: func() (slash.SlashingLevel, slash.SlashingType) {
				return slash.Level1, "invalid slash type"
			},
			expectedErr: process.ErrInvalidSlashType,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType) {
				return 44444, slash.None
			},
			expectedErr: process.ErrInvalidSlashLevel,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType) {
				return slash.Level1, slash.MultipleProposal
			},
			expectedErr: process.ErrCannotCastProofToMultipleProposedHeaders,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType) {
				return slash.Level0, slash.None
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		sd, _ := detector.NewHeaderSlashingDetector(&mock2.NodesCoordinatorStub{})
		proof := slash.NewSlashingProof(currTest.args())

		err := sd.ValidateProof(proof)
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestHeaderSlashingDetector_ValidateProof_MultipleProposalProof_InvalidSlashLevelsAndTypes_ExpectError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData)
		expectedErr error
	}{
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level0, slash.MultipleProposal, []process.InterceptedData{}
			},
			expectedErr: process.ErrInvalidSlashLevel,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return 4444, slash.MultipleProposal, []process.InterceptedData{}
			},
			expectedErr: process.ErrInvalidSlashLevel,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{}
			},
			expectedErr: process.ErrNotEnoughHeadersProvided,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(2, []byte("h1")),
					createInterceptedHeaderData(2, []byte("h2")),
					createInterceptedHeaderData(2, []byte("h3")),
				}
			},
			expectedErr: process.ErrSlashLevelDoesNotMatchSlashType,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level2, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(2, []byte("h1")),
					createInterceptedHeaderData(2, []byte("h2")),
				}
			},
			expectedErr: process.ErrSlashLevelDoesNotMatchSlashType,
		},
	}

	for _, currTest := range tests {
		sd, _ := detector.NewHeaderSlashingDetector(&mock2.NodesCoordinatorStub{})
		proof, _ := slash.NewMultipleProposalProof(currTest.args())

		err := sd.ValidateProof(proof)
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestHeaderSlashingDetector_ValidateProof_DifferentHeaders(t *testing.T) {
	t.Parallel()

	errGetProposer := errors.New("cannot get proposer")
	tests := []struct {
		args        func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData)
		expectedErr error
	}{
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(5, []byte("h1")),
					createInterceptedHeaderData(5, []byte("h1")),
				}
			},
			expectedErr: process.ErrProposedHeadersDoNotHaveDifferentHashes,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level2, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(5, []byte("h1")),
					createInterceptedHeaderData(5, []byte("h2")),
					createInterceptedHeaderData(5, []byte("h2")),
				}
			},
			expectedErr: process.ErrProposedHeadersDoNotHaveDifferentHashes,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(4, []byte("h1")),
					createInterceptedHeaderData(5, []byte("h2")),
				}
			},
			expectedErr: process.ErrHeadersDoNotHaveSameRound,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level2, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(4, []byte("h1")),
					createInterceptedHeaderData(4, []byte("h2")),
					createInterceptedHeaderData(5, []byte("h3")),
				}
			},
			expectedErr: process.ErrHeadersDoNotHaveSameRound,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(0, []byte("h1")),
					createInterceptedHeaderData(0, []byte("h2")),
				}
			},
			expectedErr: errGetProposer,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(0, []byte("h")),
					createInterceptedHeaderData(0, []byte("h1")),
				}
			},
			expectedErr: errGetProposer,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(1, []byte("h1")),
					createInterceptedHeaderData(1, []byte("h2")),
				}
			},
			expectedErr: process.ErrHeadersDoNotHaveSameProposer,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level1, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(4, []byte("h1")),
					createInterceptedHeaderData(4, []byte("h2")),
				}
			},
			expectedErr: nil,
		},
		{
			args: func() (slash.SlashingLevel, slash.SlashingType, []process.InterceptedData) {
				return slash.Level2, slash.MultipleProposal, []process.InterceptedData{
					createInterceptedHeaderData(5, []byte("h1")),
					createInterceptedHeaderData(5, []byte("h2")),
					createInterceptedHeaderData(5, []byte("h3")),
				}
			},
			expectedErr: nil,
		},
	}

	sd, _ := detector.NewHeaderSlashingDetector(
		&mock2.NodesCoordinatorStub{
			ComputeConsensusGroupCalled: func(randomness []byte, round uint64, _ uint32, _ uint32) ([]sharding.Validator, error) {
				if round == 0 && bytes.Equal(randomness, []byte("h1")) {
					return nil, errGetProposer
				}
				if round == 1 && bytes.Equal(randomness, []byte("h1")) {
					return []sharding.Validator{mock.NewValidatorMock([]byte("proposer1"))}, nil
				}
				if round == 1 && bytes.Equal(randomness, []byte("h2")) {
					return []sharding.Validator{mock.NewValidatorMock([]byte("proposer2"))}, nil
				}

				return []sharding.Validator{mock.NewValidatorMock([]byte("proposer"))}, nil
			},
		})

	for _, currTest := range tests {
		proof, _ := slash.NewMultipleProposalProof(currTest.args())
		err := sd.ValidateProof(proof)
		require.Equal(t, currTest.expectedErr, err)
	}
}

func createInterceptedHeaderArg(round uint64, randSeed []byte) *interceptedBlocks.ArgInterceptedBlockHeader {
	args := &interceptedBlocks.ArgInterceptedBlockHeader{
		ShardCoordinator:        &mock.ShardCoordinatorStub{},
		Hasher:                  &mock.HasherMock{},
		Marshalizer:             &mock.MarshalizerMock{},
		HeaderSigVerifier:       &mock.HeaderSigVerifierStub{},
		HeaderIntegrityVerifier: &mock.HeaderIntegrityVerifierStub{},
		ValidityAttester:        &mock.ValidityAttesterStub{},
		EpochStartTrigger:       &mock.EpochStartTriggerStub{},
	}

	hdr := createBlockHeaderData(round, randSeed)
	args.HdrBuff, _ = args.Marshalizer.Marshal(hdr)

	return args
}

func createBlockHeaderData(round uint64, randSeed []byte) *block.Header {
	return &block.Header{
		RandSeed: randSeed,
		ShardID:  1,
		Round:    round,
		Epoch:    3,
	}
}

func createInterceptedHeaderData(round uint64, randSeed []byte) *interceptedBlocks.InterceptedHeader {
	args := createInterceptedHeaderArg(round, randSeed)
	interceptedHeader, _ := interceptedBlocks.NewInterceptedHeader(args)

	return interceptedHeader
}
