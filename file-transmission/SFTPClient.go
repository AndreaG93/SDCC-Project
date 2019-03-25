package file_transmission

import (
	"SDCC-Project-WorkerNode/corelity"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
)

func SendFile(fullAddress string) error {

	var SSHConnection *ssh.Client
	var SFTPClient *sftp.Client
	var fileOnRemoteServer *sftp.File
	var err error

	config := ssh.ClientConfig{
		User: "Graziani",
		Auth: []ssh.AuthMethod{
			ssh.Password(""),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	if SSHConnection, err = ssh.Dial("tcp", fullAddress, &config); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(SSHConnection.Close())
	}()

	if SFTPClient, err = sftp.NewClient(SSHConnection); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(SFTPClient.Close())
	}()

	if fileOnRemoteServer, err = SFTPClient.Create("hello.txt"); err != nil {
		return err
	}

	if _, err := fileOnRemoteServer.Write([]byte("Hello world!")); err != nil {
		return err
	}

	// check it's there
	fi, err := SFTPClient.Lstat("hello.txt")
	if err != nil {
		return err
	}
	log.Println(fi)

	return nil
}
