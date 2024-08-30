package main

import (
	"fmt"

	"github.com/bsach64/goback/chunking"
)


func main() {
  f, err := chunking.ChunkFile("example.txt")
  if err!=nil{
    fmt.Println(err)
  }
  chunking.HashChunks(f) 
}


// Push this where 
