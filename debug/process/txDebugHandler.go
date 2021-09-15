package process

import (
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const maxDuration = time.Second * 10

var TX_DEBUGGER = NewTxDebugHandler()
var log = logger.GetOrCreate("debug")

type txItem struct {
	timeStamp time.Time
}

type txDebugHandler struct {
	mut     sync.RWMutex
	txMap   map[string]txItem
	history map[string]struct{}
}

func NewTxDebugHandler() *txDebugHandler {
	tdh := &txDebugHandler{
		txMap:   make(map[string]txItem),
		history: make(map[string]struct{}),
	}

	go tdh.check()

	return tdh
}

func (tdh *txDebugHandler) AddTx(hash string) {
	tdh.mut.Lock()
	defer tdh.mut.Unlock()

	_, found := tdh.history[hash]
	if found {
		return
	}

	tdh.txMap[hash] = txItem{
		timeStamp: time.Now(),
	}
}

func (tdh *txDebugHandler) RemoveTx(hash string) {
	tdh.mut.Lock()
	delete(tdh.txMap, hash)
	tdh.history[hash] = struct{}{}
	tdh.mut.Unlock()
}

func (tdh *txDebugHandler) Clear() {
	tdh.mut.Lock()
	tdh.txMap = make(map[string]txItem)
	tdh.history = make(map[string]struct{})
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
