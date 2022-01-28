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
	// try with InsecureIgnoreHostKey first
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", pS.Host), sshConfig)
	if err != nil {
		return nil, err
	}

	_, err = client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (pS *PoliteScanner) IsInterfaceNil() bool {
	return pS == nil
}
