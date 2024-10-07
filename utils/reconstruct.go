package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

func Reconstruct(snapshot Snapshot) ([]byte, error) {
	allData := make([][]byte, len(snapshot.Chunks))
	var wg sync.WaitGroup
	for _, info := range snapshot.Chunks {
		wg.Add(1)
		go func(info ChunkInfo, allData [][]byte) {
			err := readChunk(info, allData)
			if err != nil {
				log.Error("Could not read chunk %v, error:", "err", info.FileName, err)
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
	log.Info("Opening Chunk", "err", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	allData[chunkInfo.Order] = data
	return nil
}
