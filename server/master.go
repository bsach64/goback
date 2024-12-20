package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/log"

	"golang.org/x/crypto/ssh"
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
	db      utils.DBConn
}

func NewMaster(ip string) error {
	rsaPath := os.Getenv("KEY_PATH")

	if rsaPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("Cannot get the working directory:", "Error", err)
		}
		rsaPath = wd + "/private/id_rsa"
	}

	m := Server{
		index: 0,
		Host:  ip,
		Port:  2022,
		IdRsa: rsaPath,
	}

	var err error
	m.db, err = utils.CreateDBConn("meta.db")
	if err != nil {
		return err
	}

	go func() {
		err := m.ListenAndServe()
		if err != nil {
			log.Fatal("Error while creating master server", "err", err)
		}
	}()
	select {}
}

// Listen and Serve the master
func (m *Server) ListenAndServe() error {
	privateBytes, err := os.ReadFile(m.IdRsa)
	if err != nil {
		log.Error("Failed to load private key:", "Error", err)
		return err
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Error("Failed to parse private key:", "Error", err)
		return err
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(private)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen on %v", addr, "Error", err)
		return err
	}
	defer listener.Close()

	log.Info("Master SSH server listening on ", "Host", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Failed to accept connection", "Error", err)
			continue
		}

		sshConn, _, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Error("Failed to handshake,", "Error", err)
			continue
		}

		log.Info(fmt.Sprintf("New SSH to master from %s", sshConn.RemoteAddr()))
		go m.handleClient(sshConn, reqs)
	}
}

func (m *Server) handleClient(conn *ssh.ServerConn, reqs <-chan *ssh.Request) {
	for req := range reqs {
		switch req.Type {
		case "worker-details":
			var newWorker Worker
			err := json.Unmarshal(req.Payload, &newWorker)
			if err != nil {
				log.Error("could not get worker details", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte("could not get worker details"))
					if err != nil {
						log.Error("Could not send reply", "err", err)
					}
				}
				continue
			}

			err = m.db.WriteClientInfo(utils.ClientInfo{
				IP:    newWorker.Ip,
				Alive: true,
			})

			if err != nil {
				log.Error("Could not write to db", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte("Could not write to database"))
					if err != nil {
						log.Error("Could not send reply", "err", err)
					}
				}
				continue
			}

			m.addWorker(newWorker)
			if req.WantReply {
				err := req.Reply(true, []byte("Got Worker Details"))
				if err != nil {
					log.Error("could not get worker details", "err", err)
					continue
				}
			}
		case "create-backup":
			replyMessage, err := json.Marshal(m.workers)
			if err != nil {
				log.Error("failed to marshal worker node", "err", err)
				continue
			}

			log.Info("Received Create-Backup request with", "payload", string(req.Payload))
			if req.WantReply {
				err := req.Reply(true, replyMessage)
				if err != nil {
					log.Error("Cannot reply to request from", "addr", conn.RemoteAddr().String(), "err", err)
					continue
				}
			}
		case "close-connection":
			log.Info("Received close-connection request")
			// Implement logic to close the connection

			//Remove the Worker IP
			workerIP := string(req.Payload)
			replyMessage := []byte("Connection closing")
			if req.WantReply {
				err := req.Reply(true, replyMessage)
				if err != nil {
					log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
				}
			}
			log.Infof("Connection closed with %v", conn.RemoteAddr().String())
			m.RemoveWorker(workerIP)

		case "start-file-upload":
			var fileInfo utils.FileInfo
			err := json.Unmarshal(req.Payload, &fileInfo)
			if err != nil {
				log.Error("Failed to marshal file info", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte(err.Error()))
					if err != nil {
						log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
					}
				}
				continue
			}

			// Currently assuming IPv4
			// Port number of worker and client is different causing issues
			ip := conn.RemoteAddr().String()
			ip, _, _ = strings.Cut(ip, ":")
			err = m.db.StartFileUpload(ip, fileInfo.Filename, fileInfo.Size)
			if err != nil {
				log.Error("Could not write to db", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte("failed to start file upload"))
					if err != nil {
						log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
					}
				}
				continue
			}

			if req.WantReply {
				err := req.Reply(true, []byte("Saved metadata!"))
				if err != nil {
					log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
				}
			}

		case "finish-file-upload":
			var fileInfo utils.FileInfo
			err := json.Unmarshal(req.Payload, &fileInfo)
			if err != nil {
				log.Error("Failed to marshal file info", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte("failed to start file upload"))
					if err != nil {
						log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
					}
				}
				continue
			}

			// Currently assuming IPv4
			// Port number of worker and client is different causing issues
			ip := conn.RemoteAddr().String()
			ip, _, _ = strings.Cut(ip, ":")
			err = m.db.FinishFileUpload(ip, fileInfo.Filename)
			if err != nil {
				log.Error("Could not write to db", "err", err)
				if req.WantReply {
					err := req.Reply(false, []byte("failed to finish file upload"))
					if err != nil {
						log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
					}
				}
				continue
			}

			if req.WantReply {
				err := req.Reply(true, []byte("Upload Successful!"))
				if err != nil {
					log.Errorf("Cannot reply to connection from : %v", conn.RemoteAddr().String())
				}
			}

		default:
			fmt.Println("Unknown request type:", req.Type)
			err := req.Reply(false, nil) // Deny unknown requests
			if err != nil {
				log.Errorf("Replying to unknown request failed from: %v", conn.RemoteAddr().String())
			}
		}
	}
}

func (m *Server) addWorker(newWorker Worker) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.workers = append(m.workers, newWorker)
}

func (m *Server) RemoveWorker(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, worker := range m.workers {
		if worker.Ip == ip {
			m.workers = append(m.workers[:i], m.workers[i+1:]...)
			log.Info("Removed worker", "ip", ip)
			return
		}
	}
	log.Warn("Worker not found", "ip", ip)
}
