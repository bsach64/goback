package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPServer struct {
	Host  string
	Port  int
	IdRsa string
}

func New(host, id_rsa string, port int) SFTPServer {
	return SFTPServer{
		Host:  host,
		Port:  port,
		IdRsa: id_rsa,
	}
}

func Listen(s SFTPServer) {
	rsa_key := s.IdRsa
	if rsa_key == "" {
		rsa_key = "private/id_rsa"
	}
	privateBytes, err := os.ReadFile(rsa_key)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
		// TODO : Change this to false and add client authentication
	}

	config.AddHostKey(private)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on 2022: %v", err)
	}
	defer listener.Close()
	log.Println("SFTP Server listening at :", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed during incoming connection: %v", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Printf("Failed handshake: %v", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go ssh.DiscardRequests(reqs)

		for newChannel := range chans {
			if newChannel.ChannelType() != "session" {
				if newChannel.Reject(ssh.UnknownChannelType, "unknown channel type") != nil {
					log.Printf("Error while rejecting channel creation")
				}
				continue
			}

			channel, requests, err := newChannel.Accept()
			if err != nil {
				// failed channel does not kill the server
				log.Printf("Failed accepting channel: %v", err)
				continue
			}

			go func(in <-chan *ssh.Request) {
				for req := range in {
					if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
						if req.Reply(true, nil) != nil {
							log.Printf("Cannot send Reply to the request")
						}
						handleSFTP(channel)
					} else {
						if req.Reply(false, nil) != nil {
							log.Printf("Cannot send Reply to the request")
						}
					}
				}
			}(requests)
		}
	}
}

func handleSFTP(channel ssh.Channel) {
	server, err := sftp.NewServer(channel)
	if err != nil {

		// Only fatal log that exits the server creation
		log.Fatalf("Failed SFTP server creation %v", err)

		return
	}
	defer server.Close()

	if err := server.Serve(); err == io.EOF {
		log.Println("SFTP client exited")
		return
	} else if err != nil {
		log.Printf("SFTP failed with error: %v", err)
	}
}
