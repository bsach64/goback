package client

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/bsach64/goback/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	user     string
	password string
}

func NewClient(user, passwd string) Client {

	return Client{
		user:     user,
		password: passwd,
	}

}

func (c *Client) ConnectToServer(host string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            c.user,
		Auth:            []ssh.AuthMethod{ssh.Password(c.password)},
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

	file, err := utils.ChunkFile(f)

	if err != nil {
		log.Fatalf("Cannot chunk the file %v", err)
	}

	hashed_chunks := utils.HashChunks(file)

	err = uploadChunks(sftpClient, hashed_chunks)
	if err != nil {
		log.Fatalf("%v", err)
	}

}

func uploadChunks(sftpClient *sftp.Client, chunks map[string][]byte) error {
	for key, val := range chunks {
		remoteFilePath := filepath.Join("./tmp", fmt.Sprintf("%s.chunk", key)) // Writes chunks to the remote file from the byte array

		remoteFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return err
		}

		defer remoteFile.Close()

		if _, err := remoteFile.Write(val); err != nil {
			return err
		}

		log.Printf("Chunk uploaded successfully to %s", remoteFilePath)
	}
	return nil
}
