package process

import (
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const maxDuration = time.Second * 6

var TX_DEBUGGER = NewTxDebugHandler()
var log = logger.GetOrCreate("debug")

type txItem struct {
	timeStamp time.Time
}

type txDebugHandler struct {
	mut   sync.RWMutex
	txMap map[string]txItem
}

func NewTxDebugHandler() *txDebugHandler {
	tdh := &txDebugHandler{
		txMap: make(map[string]txItem),
	}

	go tdh.check()

	return tdh
}

func (tdh *txDebugHandler) AddTx(hash string) {
	tdh.mut.Lock()
	tdh.txMap[hash] = txItem{
		timeStamp: time.Now(),
	}
	tdh.mut.Unlock()
}

func (tdh *txDebugHandler) RemoveTx(hash string) {
	tdh.mut.Lock()
	delete(tdh.txMap, hash)
	tdh.mut.Unlock()
}

func (tdh *txDebugHandler) Clear() {
	tdh.mut.Lock()
	tdh.txMap = make(map[string]txItem)
	tdh.mut.Unlock()
}

func (tdh *txDebugHandler) check() {
	for {
		time.Sleep(time.Second)

		tdh.checkMap()
	}
}

func (tdh *txDebugHandler) checkMap() {
	tdh.mut.RLock()
	defer tdh.mut.RUnlock()

	for hash, tx := range tdh.txMap {
		dur := time.Since(tx.timeStamp)
		if dur > maxDuration {
			log.Error("stalled transaction found", "tx hash", []byte(hash), "stalled for", dur)
		}
	}
}
