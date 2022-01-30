package scan

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type PoliteScanner struct {
	Host string
	Port int
	User string
	Pwd  string
}

// Run polite ssh scan
func (pS *PoliteScanner) Scan() (res []byte, err error) {
	sshConfig := &ssh.ClientConfig{
		User: pS.User,
		Auth: []ssh.AuthMethod{ssh.Password(pS.Pwd)},
	}

	// InsecureIgnoreHostKey returns a function that can be used for ClientConfig.HostKeyCallback to accept any host key.
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	_, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", pS.Host, pS.Port), sshConfig)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (pS *PoliteScanner) IsInterfaceNil() bool {
	return pS == nil
}
