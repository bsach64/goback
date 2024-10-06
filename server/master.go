package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/bsach64/goback/server/backuptask"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MasterServer struct {
    workers     []WorkerNode
    index       int
    mu          sync.Mutex
    Host        string
    Port        int
    IdRsa       string
}

func NewMaster() {
   
    m := MasterServer {
        index: 0,
        Host: "0.0.0.0",
        Port: 2022,
        IdRsa: "private/id_rsa",
    }
    err := StartNewWorker(&m,1,1,"127.0.0.1",2025)
    if err!=nil{
        log.Fatalf("New Worker fuck")
    }
    go m.ListenAndServe()
    select {}
}

func (m *MasterServer) ListenAndServe() error {
    privateBytes, err := os.ReadFile(m.IdRsa)
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
    }
    config.AddHostKey(private)

    addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Printf("Failed to listen on %v: %v\n", addr, err)
        return err
    }
    defer listener.Close()

    log.Printf("Master SSH server listening on %v\n", addr)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v\n", err)
            continue
        }

        sshConn, _, reqs, err := ssh.NewServerConn(conn, config)
        if err != nil {
            log.Printf("Failed to handshake: %v\n", err)
            continue
        }

        log.Printf("New SSH connection from %s\n", sshConn.RemoteAddr())
        go m.handleClient(sshConn,reqs)
    }
}

func (m *MasterServer) handleClient(conn *ssh.ServerConn, reqs <-chan *ssh.Request) {
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
                    req.Reply(true, replyMessage)
                }

            case "list-backups":
                fmt.Println("Received list-backups request")
                backupList := []byte("Backup1, Backup2, Backup3")
                if req.WantReply {
                    req.Reply(true, backupList)
                }

            case "close-connection":
                fmt.Println("Received close-connection request")
                // Implement logic to close the connection
                replyMessage := []byte("Connection closing")
                if req.WantReply {
                    req.Reply(true, replyMessage)
                }
                conn.Close()

            default:
                fmt.Println("Unknown request type:", req.Type)
                req.Reply(false, nil) // Deny unknown requests
        }
    }   
}

func (m *MasterServer) assignWorker(fileName string, fileSize int64) (*pb.WorkerAssignmentResponse, error) {
    worker := m.chooseWorker() 

    conn, err := grpc.NewClient(worker.gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := pb.NewMasterServiceClient(conn)

    req := &pb.BackupTaskRequest{
        FileName: fileName,
        FileSize: fileSize,
    }

    resp, err := client.RequestWorker(context.Background(), req)
    if err != nil {
        return nil, err
    }

    return resp, nil
}


func (m *MasterServer) chooseWorker() WorkerNode {
  	m.mu.Lock()
	defer m.mu.Unlock()

	selectedWorker := m.workers[m.index]
	m.index = (m.index + 1) % len(m.workers) // Use simple Round Robin to select slave node

	return selectedWorker
}
