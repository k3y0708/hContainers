package util

import (
	"bytes"
	"k3y0708/hContainers/global"
	"net"

	"golang.org/x/crypto/ssh"
)

func RemoteRun(addr string, cmd string) (string, error) {
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
