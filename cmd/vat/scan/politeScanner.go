package scan

import "sync"

type PoliteScanner struct {
	mutScanner sync.Mutex
	Host       string
}

// Run polite ssh scan
func (pS *PoliteScanner) Scan() (res []byte, err error) {
	pS.mutScanner.Lock()
	defer pS.mutScanner.Unlock()

	dialer := newSshDialer(pS.Host)
	err = dialer.Dial()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (pS *PoliteScanner) IsInterfaceNil() bool {
	return pS == nil
}
