package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/charmbracelet/log"

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

func (s *SFTPServer) Listen() error {
	rsaKey := s.IdRsa
	if rsaKey == "" {
		rsaKey = "private/id_rsa"
	}
	privateBytes, err := os.ReadFile(rsaKey)
	if err != nil {
		log.Info("Failed to load private key:", "err", err)
		return err
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Info("Failed to parse private key:", "err", err)
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
		log.Info("Failed to listen on %v:", "err", s.Port, err)
		return err
	}
	defer listener.Close()
	log.Info("Worker Server listening at", "IP",addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Info("Failed during incoming connection:", "err", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Info("Failed handshake:", "err", err)
			continue
		}

		log.Printf("SSH connection request %s - %s ", sshConn.RemoteAddr(), addr)
		go ssh.DiscardRequests(reqs)

		for newChannel := range chans {
			if newChannel.ChannelType() != "session" {
				if newChannel.Reject(ssh.UnknownChannelType, "unknown channel type") != nil {
					log.Info("Error while rejecting channel creation\n")
				}
				continue
			}

			channel, requests, err := newChannel.Accept()
			if err != nil {
				// failed channel does not kill the server
				log.Info("Failed accepting channel:", "err", err)
				continue
			}

			go func(in <-chan *ssh.Request) {
				for req := range in {
					if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
						if req.Reply(true, nil) != nil {
							log.Info("Cannot send Reply to the request\n")
						}
						err = handleSFTP(channel)
						if err != nil {
							log.Info("Could not handle SFTP Connection", "err", err)
						}
					} else {
						if req.Reply(false, nil) != nil {
							log.Info("Cannot send Reply to the request\n")
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
		log.Info("SFTP client exited")
	} else if err != nil {
		return err
	}
	return nil
}
