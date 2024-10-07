package client

import (
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"github.com/bsach64/goback/server"
)

func TestClient(t *testing.T) {
	s := server.New("0.0.0.1", "../private/id_rsa", 2022)
	go func(s server.SFTPServer) {
		err := s.Listen()
		if err != nil {
			log.Error("Could not start server:", "err", err)
		}
	}(s)
	go t.Run("Connection Test", func(t *testing.T) {
		testConnection(t)
	})
	go t.Run("Upload Test", func(t *testing.T) {
		testUpload(t)
	})
}

func testConnection(t *testing.T) {
	client := NewClient("test_user", "test_password")

	_, err := client.ConnectToServer("127.0.0.1:2022")

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Refused Connection from Server")
			t.FailNow()
		} else {
			t.Errorf("%v", err)
		}
		return
	}

}

func testUpload(t *testing.T) {

	client := NewClient("test_user", "test_password")

	ssh_client, err := client.ConnectToServer("127.0.0.1:2022")

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Refused Connection from Server")
		} else {
			t.Errorf("%v", err)
		}
		return
	}

	err = Upload(ssh_client, "../test_files/example.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

}
