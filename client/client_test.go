package client

import (
	"strings"
	"testing"

	"github.com/bsach64/goback/server"
)

func TestClient(t *testing.T){
  go server.Listen("../private/id_rsa");
  go t.Run("Connection Test",func (t *testing.T)  {
    testConnection(t)
  })
  go t.Run("Upload Test",func (t *testing.T)  {
    testUpload(t)
  })
}

func testConnection(t *testing.T) {

    _, err := ConnectToServer(
      "test_user",
      "password",
      "127.0.0.1:2022",
    )

     if err!=nil{
      if strings.Contains(err.Error(),"connection refused"){
        t.Skipf("Refused Connection from Server")
        t.FailNow()
      }else{
        t.Errorf("%v",err)
      }
      return
   } 

}

func testUpload(t *testing.T) {

    client, err := ConnectToServer(
      "test_user",
      "password",
      "127.0.0.1:2022",
    )


    if err!=nil{
      if strings.Contains(err.Error(),"connection refused"){
        t.Skipf("Refused Connection from Server")
      }else{
        t.Errorf("%v",err)
      }
      return
   } 

    err = Upload(client,"../test_files/example.txt")
    if err!=nil{
      t.Errorf("Test Upload failed : %v",err)
      return
    }
  
}

