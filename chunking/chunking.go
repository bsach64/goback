package chunking

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aclements/go-rabin/rabin"
)

// Chunked File --> Many Chunked Files --> Pack
// Pack is finally stored in the backup device


/// Use case for meta data is :
/// Allows the design of the metadata 
/// to evolve without re-uploading the blobs.

type MetaData struct{
  processed_at time.Time
  file_name string
}

// File is broken into its meta data and chunks 
type File struct{
  meta MetaData;
  file [][]byte
}

func ChunkFile(filename string) (File,error) {
    var result File
    file, err := os.Open(filename)
    if err != nil {
        return result,fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    const minSize = 512     // Minimum chunk size in bytes
    const avgSize = 2048    // Average chunk size in bytes
    const maxSize = 8192    // Maximum chunk size in bytes

    var chunk_buffer [][]byte
    chunker := rabin.NewChunker(rabin.NewTable(rabin.Poly64,256),file,minSize,avgSize,maxSize);
    // Window size <= minSize (256 in this case)-------------^

    for {
      chunk, err := chunker.Next() 
      if err == io.EOF {
        break 
      }
      if err != nil {
        return result,fmt.Errorf("error reading chunk: %v", err)
      }
      // instead we should return the slices of file into a buffer
      buffer := make([]byte, chunk)
      file.Read(buffer)
      chunk_buffer = append(chunk_buffer,buffer)
      fmt.Println(string(buffer))  
  }
    result.file = chunk_buffer

    result.meta.file_name = filename
    result.meta.processed_at = time.Now() 
    return result,nil
}


