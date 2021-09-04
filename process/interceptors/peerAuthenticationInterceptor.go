package interceptors

import (
	"errors"
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/debug/resolver"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
)

const blackListCauseWrongTypeAssertion = "intercepted data is not of type process.InterceptedPeerAuthentication"

// ArgPeerAuthenticationInterceptor is the argument for the peer authentication interceptor
type ArgPeerAuthenticationInterceptor struct {
	ArgSingleDataInterceptor
	Marshalizer             marshal.Marshalizer
	ValidatorChecker        process.ValidatorChecker
	AuthenticationProcessor process.PeerAuthenticationProcessor
	ObserversThrottler      process.InterceptorThrottler
}

type peerAuthenticationInterceptor struct {
	*baseDataInterceptor
	validatorChecker            process.ValidatorChecker
	peerAuthenticationProcessor process.PeerAuthenticationProcessor
	observersThrottler          process.InterceptorThrottler
	whiteListRequest            process.WhiteListHandler
}

// NewPeerAuthenticationInterceptor hooks a new interceptor for packed multi data containing peer authentication instances
func NewPeerAuthenticationInterceptor(arg ArgPeerAuthenticationInterceptor) (*peerAuthenticationInterceptor, error) {
	err := checkArguments(arg.ArgSingleDataInterceptor)
	if err != nil {
		return nil, err
	}
	if check.IfNil(arg.Marshalizer) {
		return nil, process.ErrNilMarshalizer
	}
	if check.IfNil(arg.ValidatorChecker) {
		return nil, process.ErrNilValidatorChecker
	}
	if check.IfNil(arg.AuthenticationProcessor) {
		return nil, process.ErrNilAuthenticationProcessor
	}
	if check.IfNil(arg.ObserversThrottler) {
		return nil, fmt.Errorf("%w for the observers throttler", process.ErrNilInterceptorThrottler)
	}

	interceptor := &peerAuthenticationInterceptor{
		baseDataInterceptor: &baseDataInterceptor{
			throttler:            arg.Throttler,
			antifloodHandler:     arg.AntifloodHandler,
			topic:                arg.Topic,
			currentPeerId:        arg.CurrentPeerId,
			processor:            arg.Processor,
			debugHandler:         resolver.NewDisabledInterceptorResolver(),
			marshalizer:          arg.Marshalizer,
			factory:              arg.DataFactory,
			preferredPeersHolder: arg.PreferredPeersHolder,
		},
		validatorChecker:            arg.ValidatorChecker,
		peerAuthenticationProcessor: arg.AuthenticationProcessor,
		observersThrottler:          arg.ObserversThrottler,
		whiteListRequest:            arg.WhiteListRequest,
	}

	return interceptor, nil
}

// ProcessReceivedMessage is the callback func from the p2p.Messenger and will be called each time a new message was received
// (for the topic this validator was registered to)
func (pai *peerAuthenticationInterceptor) ProcessReceivedMessage(message p2p.MessageP2P, fromConnectedPeer core.PeerID) error {
	multiDataBuff, _, err := pai.preProcessMessage(message, fromConnectedPeer)
	if err != nil {
		return err
	}

	authMessageIgnored := false
	for _, dataBuff := range multiDataBuff {
		var interceptedData process.InterceptedData
		interceptedData, err = pai.interceptedData(dataBuff, message.Peer(), fromConnectedPeer)
		if err != nil {
			pai.throttler.EndProcessing()
			return err
		}

		peerAuth, ok := interceptedData.(process.InterceptedPeerAuthentication)
		if !ok {
			//intercepted data is not of type InterceptedPeerAuthentication
			pai.blackListPeers(blackListCauseWrongTypeAssertion, nil, message.Peer(), fromConnectedPeer)
			pai.throttler.EndProcessing()

			return errors.New(blackListCauseWrongTypeAssertion)
		}

		if !pai.shouldProcess(message, peerAuth) {
			authMessageIgnored = true
			continue
		}

		pai.observersThrottler.StartProcessing()
		errProcess := pai.peerAuthenticationProcessor.ProcessReceived(dataBuff, peerAuth)
		if errProcess != nil {
			pai.throttler.EndProcessing()
			pai.observersThrottler.EndProcessing()
			pai.blackListPeers("peer info processing error", errProcess, message.Peer(), fromConnectedPeer)
			return errProcess
		}
		pai.observersThrottler.EndProcessing()
	}
	pai.throttler.EndProcessing()
	shouldNotPropagateMessage := authMessageIgnored || len(multiDataBuff) > 1
	if shouldNotPropagateMessage {
		return process.ErrShouldNotBroadcastMessage
	}

	return nil
}

func (pai *peerAuthenticationInterceptor) shouldProcess(message p2p.MessageP2P, peerAuth process.InterceptedPeerAuthentication) bool {
	isFromOtherPeer := peerAuth.PeerID() != message.Peer()
	isUnrequested := isFromOtherPeer && !pai.whiteListRequest.IsWhiteListed(peerAuth)
	if isUnrequested {
		// we do not allow propagation of peer auth messages sent form other peers other than self except
		// when we specifically do a request
		return false
	}

	_, shardID, err := pai.validatorChecker.GetValidatorWithPublicKey(peerAuth.PublicKey())
	peerAuth.SetComputedShardID(shardID)

	isValidator := err == nil
	isHardforkMessage := pai.peerAuthenticationProcessor.IsHardforkMessage(peerAuth)
	messageIsImportant := isValidator || isHardforkMessage
	if messageIsImportant {
		return true
	}

	return pai.observersThrottler.CanProcess()
}

// IsInterfaceNil returns true if there is no value under the interface
func (pai *peerAuthenticationInterceptor) IsInterfaceNil() bool {
	return pai == nil
}
