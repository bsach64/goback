package utils

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

func Reconstruct(snapshot Snapshot) ([]byte, error) {
	allData := make([][]byte, len(snapshot.Chunks))
	var wg sync.WaitGroup
	for _, info := range snapshot.Chunks {
		wg.Add(1)
		go func(info ChunkInfo, allData [][]byte) {
			err := readChunk(info, allData)
			if err != nil {
				log.Fatalf("Could not read chunk %v, error: %v\n", info.FileName, err)
			}
			wg.Done()
		}(info, allData)
	}
	wg.Wait()
	concatData := make([]byte, 0)
	for _, s := range allData {
		concatData = append(concatData, s...)
	}
	return concatData, nil
}

func readChunk(chunkInfo ChunkInfo, allData [][]byte) error {
	filePath := filepath.Join("./.data", chunkInfo.FileName)
	log.Printf("Opening Chunk %v\n", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	allData[chunkInfo.Order] = data
	return nil
}
