package utils

import (
	"log"
	"testing"

	"github.com/bsach64/goback/server"
)

func TestReconstruct(t *testing.T) {

	s := server.New("0.0.0.1", "../private/id_rsa", 2022)
	go func(s server.SFTPServer) {
		err := server.Listen(s)
		if err != nil {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}(s)
	go t.Run("Reconstruction Test:", func(t *testing.T) {

	})
}
