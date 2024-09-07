package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bsach64/goback/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	user      string
	password  string
	SSHClient *ssh.Client
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

func Upload(client *ssh.Client, f string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalf("Failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	file, err := utils.ChunkFile(f)

	if err != nil {
		return fmt.Errorf("Cannot chunk the file %v because of %v\n", f, err)
	}

	hashedChunks := utils.HashChunks(file)

	err = uploadChunks(sftpClient, hashedChunks)
	if err != nil {
		return err
	}

	err = createSnapshot(sftpClient, file, hashedChunks)
	if err != nil {
		return err
	}
	return nil
}

func uploadChunks(sftpClient *sftp.Client, chunks map[string]utils.Chunk) error {
	err := sftpClient.MkdirAll("./.data")
	if err != nil {
		return err
	}
	for key, val := range chunks {
		remoteFilePath := filepath.Join("./.data", fmt.Sprintf("%s.chunk", key)) // Writes chunks to the remote file from the byte array

		remoteFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return err
		}

		defer remoteFile.Close()

		// Can be condensed into one write call
		if _, err := remoteFile.Write([]byte(fmt.Sprintf("%v\n", val.Order))); err != nil {
			return err
		}

		if _, err := remoteFile.Write(val.Data); err != nil {
			return err
		}

		log.Printf("Chunk uploaded successfully to %s", remoteFilePath)
	}
	return nil
}

func createSnapshot(sftpClient *sftp.Client, file utils.File, chunks map[string]utils.Chunk) error {
	err := sftpClient.MkdirAll("./.data")
	if err != nil {
		return err
	}
	CleanedFileName := strings.ReplaceAll(file.Meta.FileName, "/", "-")
	remoteFilePath := filepath.Join("./.data", fmt.Sprintf("(%s)-%v.snapshot", CleanedFileName, file.Meta.ProcessedAt.Unix()))

	remoteFile, err := sftpClient.OpenFile(remoteFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	for key := range chunks {
		chunkFileName := fmt.Sprintf("%s.chunk\n", key)
		if _, err := remoteFile.Write([]byte(chunkFileName)); err != nil {
			return err
		}
	}

	log.Printf("Created snapshot for %v.\n", file.Meta.FileName)
	return nil
}
