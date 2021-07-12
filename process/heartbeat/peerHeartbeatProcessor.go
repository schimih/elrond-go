package heartbeat

import (
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/syndtr/goleveldb/leveldb/cache"
)

type peerHeartbeatProcessor struct {
	cache                    cache.Cacher
	peerSignatureHandler     crypto.PeerSignatureHandler
	marshalizer              marshal.Marshalizer
	networkShardingCollector process.NetworkShardingCollector
}

// ProcessReceived processes the received message
func (php *peerHeartbeatProcessor) ProcessReceived(message p2p.MessageP2P, peerHeartbeat process.InterceptedPeerAuthentication) error {

	err := php.peerSignatureHandler.VerifyPeerSignature(peerHeartbeat.PublicKey(), peerHeartbeat.PeerID(), peerHeartbeat.Signature())
	if err != nil {
		return err
	}

	php.networkShardingCollector.UpdatePeerIdPublicKey(peerHeartbeat.PeerID(), peerHeartbeat.PublicKey())
	//add into the last failsafe map. Useful for observers.
	php.networkShardingCollector.UpdatePeerIdShardId(message.Peer(), peerHeartbeat.ComputedShardID())

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (php *peerHeartbeatProcessor) IsInterfaceNil() bool {
	return php == nil
}
