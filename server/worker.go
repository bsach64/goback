package server

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pb "github.com/bsach64/goback/server/backuptask"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Worker is just a usual SFTP server that handles the file request
// Master assigns client with a worker server
type Worker struct {
	id          int
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	SftpUser    string `json:"sftpUser"`
	SftpPass    string `json:"sftpPass"`
	sftpServer  *SFTPServer
	master      *Server
	gRPCAddress string
}

func StartNewWorker(m *Server, id int, grpcPort int, workerIP string, port int) error {

	w := Worker{
		master:      m,
		id:          id,
		Ip:          workerIP,
		Port:        port,
		SftpUser:    "sftpUser",
		SftpPass:    "sftpPass",
		gRPCAddress: fmt.Sprintf("%s:%d", workerIP, grpcPort), // Set the gRPC address
	}
	m.mu.Lock()

	defer m.mu.Unlock()

	m.workers = append(m.workers, w)

	w.StartSFTPServer()

	return nil

}

// gRPC to be added for ->
// 1. Master will not store the Credentials of the worker it will ask for it using gRPC
// 2. When Worker finishes the job the Worker informs the master using gRPC
// 3. Heartbeat checkup (worker responds with the status of being available or not)
type WorkerService struct {
	pb.UnimplementedMasterServiceServer
	worker *Worker
}

func (w *WorkerService) RequestWorker(ctx context.Context, req *pb.BackupTaskRequest) (*pb.WorkerAssignmentResponse, error) {
	log.Printf("Worker %d received task: %s\n", w.worker.id, req.FileName)

	return &pb.WorkerAssignmentResponse{
		WorkerIp:     w.worker.Ip,
		WorkerPort:   int32(w.worker.Port),
		SftpUsername: w.worker.SftpUser,
		SftpPassword: w.worker.SftpPass,
	}, nil
}

func (w *WorkerService) ReportWorkerStatus(ctx context.Context, req *pb.WorkerStatusRequest) (*pb.WorkerStatusResponse, error) {
	log.Printf("Worker %d status requested\n", w.worker.id)

	// Process worker status, e.g., is worker available
	if req.IsAvailable {
		return &pb.WorkerStatusResponse{
			Status: "Worker is available",
		}, nil
	}

	return &pb.WorkerStatusResponse{
		Status: "Worker is not available",
	}, nil
}

func (w *Worker) StartSFTPServer() {
	sftpServer := New(w.Ip, "private/id_rsa", w.Port)
	w.sftpServer = &sftpServer

	go func() {
		err := w.sftpServer.Listen()
		if err != nil {
			log.Fatalf("Worker SFTP server failed to listen: %v", err)
		}
	}()
}

func (w *Worker) ReportTaskCompletion() error {

	masterAddress := fmt.Sprintf("%s:%d", w.master.Host, w.master.Port)

	conn, err := grpc.NewClient(masterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewMasterServiceClient(conn)

	req := &pb.WorkerStatusRequest{
		WorkerId:    strconv.Itoa(w.id),
		IsAvailable: true,
	}

	_, err = client.ReportWorkerStatus(context.Background(), req)
	return err
}
