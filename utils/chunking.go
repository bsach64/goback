package utils

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"time"

	"github.com/aclements/go-rabin/rabin"
)

// File --> Many Chunks --> Many Chunks form a pack
// Pack is finally stored in the backup device

/// Use case for meta data is :
/// Allows the design of the metadata
/// to evolve without re-uploading the blobs.

const minSize = 512  // Minimum chunk size in bytes
const avgSize = 2048 // Average chunk size in bytes
const maxSize = 8192 // Maximum chunk size in bytes

type MetaData struct {
	processed_at time.Time
	file_name    string
	size         int64
}

// File is broken into its meta data and chunks
type File struct {
	meta MetaData
	file [][]byte
}

func ChunkFile(filename string) (File, error) {
	var result File
	file, err := os.Open(filename)
	if err != nil {
		return result, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var chunk_buffer [][]byte
	size := 0
	chunker := rabin.NewChunker(rabin.NewTable(rabin.Poly64, 256), file, minSize, avgSize, maxSize)
	// Window size <= minSize (256 in this case)-------------^
	for {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, fmt.Errorf("error reading chunk: %v", err)
		}
		// instead we should return the slices of file into a buffer
		buffer := make([]byte, chunk)
		_, err = file.ReadAt(buffer, int64(size))
		if err != nil {
			log.Printf("Error while reading the bytes of the file %v", err)
		}
		size += chunk
		chunk_buffer = append(chunk_buffer, buffer)
	}
	result.file = chunk_buffer

	result.meta.file_name = filename
	result.meta.processed_at = time.Now()
	result.meta.size = int64(size)
	return result, nil
}

// Hash each chunk using fnv64 as its fast and low collison
// Hashing will avoid chunk duplication

func HashChunks(f File) (map[string][]byte, error) {
	chunks := f.file
	hash_map := make(map[string][]byte)
	hash := fnv.New64()

	for _, chunk := range chunks {
		hash.Write(chunk)
		hashStr := fmt.Sprintf("%x", hash.Sum64())
		//fmt.Printf("Len of chunk : %d :: Hash: %s\n",len(chunk),hashStr)
		hash_map[hashStr] = chunk
		hash.Reset()
	}
	return hash_map, nil
}
