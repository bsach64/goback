package client

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func ConnectToServer(user, password, host string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Upload(client *ssh.Client, f string) {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalf("Failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	localFile, err := os.Open(f)
	if err != nil {
		log.Fatalf("Failed to open local file: %v", err)
	}
	defer localFile.Close()

	// Ensure the remote directory structure exists
	remoteDir := "./tmp"
	err = sftpClient.MkdirAll(remoteDir)
	if err != nil {
		log.Fatalf("Failed to create remote directory structure: %v", err)
	}

	// Create the file on the remote server
	remoteFilePath := filepath.Join(remoteDir, filepath.Base(f))
	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		log.Fatalf("Failed to create remote file: %v", err)
	}
	defer remoteFile.Close()

	// Copy the local file to the remote file
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}

	log.Printf("File uploaded successfully to %s", remoteFilePath)
}
