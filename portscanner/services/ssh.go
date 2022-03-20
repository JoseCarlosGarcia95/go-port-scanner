package services

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// IsSSH checks if a ssh service is running on a given port.
func IsSSH(host string, port uint32) bool {
	resultChain := make(chan bool, 1)

	go func() {
		sshConfig := &ssh.ClientConfig{
			User:    "portscanner",
			Auth:    []ssh.AuthMethod{ssh.Password("portscanner")},
			Timeout: time.Second,
		}

		sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

		client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(int(port)), sshConfig)

		if err != nil {
			if strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate") {
				resultChain <- true
				return
			}

			return
		}
		defer client.Close()

		return
	}()

	timeout := time.Second
	select {
	case result := <-resultChain:
		close(resultChain)
		return result
	case <-time.After(timeout):
		return false
	}
}
