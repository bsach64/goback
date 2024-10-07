package server

import (
	// "context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"net"
	"os"
	"sync"

	// pb "github.com/bsach64/goback/server/backuptask"
	"golang.org/x/crypto/ssh"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

// Master server is defined using
// 1. List of workers
// 2. Index of last worker that was assigned task
// 3. Host host for the Master currently (0.0.0.0)
// 4. Port of the Master
// 5. IdRsa since this is a SSH server the master should also have a IDRsa
type Server struct {
	workers []Worker
	index   int
	mu      sync.Mutex
	Host    string
	Port    int
	IdRsa   string
}

// Creates a new master server at 0.0.0.0 at port 2022
func NewMaster() {

	//Master server
	m := Server{
		index: 0,
		Host:  "0.0.0.0",
		Port:  2022,
		IdRsa: "private/id_rsa",
	}

	//Create a new worker to be changed later
	err := StartNewWorker(&m, 1, 1, "127.0.0.1", 2025)
	if err != nil {
		log.Fatalf("Creation of new worker failed")
	}
	go func() {
		err := m.ListenAndServe()
		if err != nil {
			log.Fatalf("Error while creating master")
		}
	}()
	select {}
}

// Listen and Serve the master
func (m *Server) ListenAndServe() error {
	privateBytes, err := os.ReadFile(m.IdRsa)
	if err != nil {
		log.Error("Failed to load private key:", "Error",err)
		return err
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Error("Failed to parse private key:",  "Error",err)
		return err
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(private)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen on %v", addr,"Error", err)
		return err
	}
	defer listener.Close()

	log.Info("Master SSH server listening on ", "Host",addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Failed to accept connection","Error", err)
			continue
		}

		sshConn, _, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Error("Failed to handshake,","Error", err)
			continue
		}

		log.Printf("New SSH to master from %s\n", sshConn.RemoteAddr())
		go m.handleClient(sshConn, reqs)
	}
}

func (m *Server) handleClient(conn *ssh.ServerConn, reqs <-chan *ssh.Request) {
	for req := range reqs {
		switch req.Type {
		case "create-backup":
			worker := m.chooseWorker()
			replyMessage, err := json.Marshal(worker)
			if err != nil {
				log.Fatalf("failed to marshal worker node: %v", err)
			}
			fmt.Println("Received create-backup request with payload:", string(req.Payload))
			if req.WantReply {

				err := req.Reply(true, replyMessage)
				if err != nil {
					fmt.Println("Cannot reply to request from :", conn.RemoteAddr().String())
				}
			}

		case "list-backups":
			fmt.Println("Received list-backups request")
			backupList := []byte("Backup1, Backup2, Backup3")
			if req.WantReply {
				err := req.Reply(true, backupList)
				if err != nil {
					fmt.Println("Cannot reply to request from :", conn.RemoteAddr().String())
				}
			}

		case "close-connection":
			fmt.Println("Received close-connection request")
			// Implement logic to close the connection
			replyMessage := []byte("Connection closing")
			if req.WantReply {
				err := req.Reply(true, replyMessage)
				if err != nil {
					fmt.Println("Cannot close connection from :", conn.RemoteAddr().String())
				}
			}
			conn.Close()

		default:
			fmt.Println("Unknown request type:", req.Type)
			err := req.Reply(false, nil) // Deny unknown requests
			if err != nil {
				fmt.Println("Replying to unknown request failed from:", conn.RemoteAddr().String())
			}
		}
	}
}

// func (m *Server) assignWorker(fileName string, fileSize int64) (*pb.WorkerAssignmentResponse, error) {
// 	worker := m.chooseWorker()

// 	conn, err := grpc.NewClient(worker.gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	client := pb.NewMasterServiceClient(conn)

// 	req := &pb.BackupTaskRequest{
// 		FileName: fileName,
// 		FileSize: fileSize,
// 	}

// 	resp, err := client.RequestWorker(context.Background(), req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

func (m *Server) chooseWorker() Worker {
	m.mu.Lock()
	defer m.mu.Unlock()

	selectedWorker := m.workers[m.index]
	m.index = (m.index + 1) % len(m.workers) // Use simple Round Robin to select slave node

	return selectedWorker
}
