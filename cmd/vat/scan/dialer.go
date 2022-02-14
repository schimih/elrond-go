package scan

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type sshDialer struct {
	host     string
	port     int
	user     string
	password string
}

func newSshDialer(host string) *sshDialer {
	return &sshDialer{
		host:     host,
		port:     DEFAULT_SSH_PORT,
		user:     "testUsername",
		password: "testPassword",
	}
}

func (dialer *sshDialer) Dial() (err error) {
	sshConfig := &ssh.ClientConfig{
		User: dialer.user,
		Auth: []ssh.AuthMethod{ssh.Password(dialer.password)},
	}

	// InsecureIgnoreHostKey returns a function that can be used for ClientConfig.HostKeyCallback to accept any host key.
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	_, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", dialer.host, dialer.port), sshConfig)
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sC *sshDialer) IsInterfaceNil() bool {
	return sC == nil
}
