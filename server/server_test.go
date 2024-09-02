package server

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestKey(t *testing.T) {
	s := New("0.0.0.1", "../private/id_rsa", 2022)
	_, err := os.Stat(s.IdRsa)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist:", s.IdRsa)
		t.Fatalf("ID_RSA file does not exist")
	} else {
		fmt.Println("File exists:", s.IdRsa)
	}
}


// Test the listen function
func TestListen(t *testing.T){
	s := New("0.0.0.0", "../private/id_rsa", 2022)
    go Listen(s)
    time.Sleep(1*time.Second)
    conn, err := net.Dial("tcp", "127.0.0.1:2022")
    if err != nil {
		t.Fatalf("Failed to connect to the SFTP server: %v", err)
	}
	conn.Close()

	t.Log("SFTP server started and accepted connections successfully")

    
}