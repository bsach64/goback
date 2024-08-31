package server

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)
func Listen() {
    privateBytes, err := os.ReadFile("private/id_rsa")
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

    listener, err := net.Listen("tcp", "0.0.0.0:2022")
    if err != nil {
        log.Fatalf("Failed to listen on 2022: %v", err)
    }
    defer listener.Close()
    log.Println("SFTP Server listening at 0.0.0.0:2022")

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
                newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
                continue
            }

            channel, requests, err := newChannel.Accept()
            if err != nil {
                log.Printf("Failed accepting channel: %v", err)
                continue
            }

            go func(in <-chan *ssh.Request) {
                for req := range in {
                    if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
                        req.Reply(true, nil)
                        handleSFTP(channel)
                    } else {
                        req.Reply(false, nil)
                    }
                }
            }(requests)
        }
    }
}

func handleSFTP(channel ssh.Channel) {
    server, err := sftp.NewServer(channel)
    if err != nil {
        log.Printf("Failed SFTP server creation %v", err)
        return
    }
    defer server.Close()

    if err := server.Serve(); err == io.EOF {
        log.Println("SFTP client exited")
    } else if err != nil {
        log.Printf("SFTP failed with error: %v", err)
    }
}


