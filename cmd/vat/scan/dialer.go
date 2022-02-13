package scan

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type sshConfig struct {
	host     string
	port     int
	user     string
	password string
}

func newDialer(host string) Dialer {
	return &sshConfig{
		host:     host,
		port:     22,
		user:     "testUsername",
		password: "testPassword",
	}
}

func (sC *sshConfig) Dial() (err error) {
	sshConfig := &ssh.ClientConfig{
		User: sC.user,
		Auth: []ssh.AuthMethod{ssh.Password(sC.password)},
	}

	// InsecureIgnoreHostKey returns a function that can be used for ClientConfig.HostKeyCallback to accept any host key.
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	_, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sC.host, sC.port), sshConfig)
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sC *sshConfig) IsInterfaceNil() bool {
	return sC == nil
}
