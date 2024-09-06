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

func New(host, idRsa string, port int) SFTPServer {
	return SFTPServer{
		Host:  host,
		Port:  port,
		IdRsa: idRsa,
	}
}

func Listen(s SFTPServer) error {
	rsaKey := s.IdRsa
	if rsaKey == "" {
		rsaKey = "private/id_rsa"
	}
	privateBytes, err := os.ReadFile(rsaKey)
	if err != nil {
		log.Printf("Failed to load private key: %v\n", err)
		return err
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Printf("Failed to parse private key: %v\n", err)
		return err
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
		// TODO : Change this to false and add client authentication
	}

	config.AddHostKey(private)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to listen on %v: %v\n", s.Port, err)
		return err
	}
	defer listener.Close()
	log.Println("SFTP Server listening at :", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed during incoming connection: %v\n", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Printf("Failed handshake: %v\n", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go ssh.DiscardRequests(reqs)

		for newChannel := range chans {
			if newChannel.ChannelType() != "session" {
				if newChannel.Reject(ssh.UnknownChannelType, "unknown channel type") != nil {
					log.Printf("Error while rejecting channel creation\n")
				}
				continue
			}

			channel, requests, err := newChannel.Accept()
			if err != nil {
				// failed channel does not kill the server
				log.Printf("Failed accepting channel: %v\n", err)
				continue
			}

			go func(in <-chan *ssh.Request) {
				for req := range in {
					if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
						if req.Reply(true, nil) != nil {
							log.Printf("Cannot send Reply to the request\n")
						}
						err = handleSFTP(channel)
						if err != nil {
							log.Printf("Could not handle SFTP Connection %v\n", err)
						}
					} else {
						if req.Reply(false, nil) != nil {
							log.Printf("Cannot send Reply to the request\n")
						}
					}
				}
			}(requests)
		}
	}
}

func handleSFTP(channel ssh.Channel) error {
	server, err := sftp.NewServer(channel)
	if err != nil {
		return err
	}
	defer server.Close()

	if err := server.Serve(); err == io.EOF {
		log.Println("SFTP client exited")
	} else if err != nil {
		return err
	}
	return nil
}
