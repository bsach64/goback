package chunking;

import (
	"os"
	"testing"

)

func TestSize(t *testing.T){
  file := "../example.jpg"

  img,err := os.Open(file)
  if err!=nil{
    t.Skipf("Image File not found %v",err)
  }

  stat,stat_err := img.Stat()
  
  if stat_err!=nil{
    t.Skipf("Cannot get stat for image file")
  }
  chunked_file,err := ChunkFile(file)
  if err!=nil{
    t.Fatalf("Chunking failed")
  }

  if chunked_file.meta.size != stat.Size(){
    t.Fatalf("Size not same Actual Size: %d Chunked total size: %d",stat.Size(),chunked_file.meta.size)
  }

}

