package util

import (
	"bytes"
	"net"
	"time"

	"github.com/hContainers/hContainers/global"

	"golang.org/x/crypto/ssh"
)

func remoteRunWithoutRetry(addr string, cmd string) (string, error) {
	// privateKey could be read from a file, or retrieved from another storage
	// source, such as the Secret Service / GNOME Keyring
	key, err := ssh.ParsePrivateKey([]byte(global.PrivateKey))
	if err != nil {
		return "", err
	}
	// Authentication
	config := &ssh.ClientConfig{
		User:            "hContainers",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	// Connect
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "2222"), config)
	if err != nil {
		return "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer  // import "bytes"
	session.Stdout = &b // get output

	// Finally, run the command
	err = session.Run(cmd)
	return b.String(), err
}

func RemoteRun(addr string, cmd string, retry int) (string, error) {
	var err error
	var output string
	broken := false
	for i := 0; i <= retry; i++ {
		output, err = remoteRunWithoutRetry(addr, cmd)
		if err == nil {
			broken = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	if !broken {
		return output, err
	} else {
		return output, nil
	}
}
