package utils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func Reconstruct(snapshotFile *os.File) ([]byte, error) {
	chunkFileNames := make([]string, 0)
	scanner := bufio.NewScanner(snapshotFile)
	for scanner.Scan() {
		nameBytes := scanner.Bytes()
		chunkFileNames = append(chunkFileNames, string(nameBytes))
	}
	allData := make([][]byte, len(chunkFileNames))
	var wg sync.WaitGroup
	for _, name := range chunkFileNames {
		wg.Add(1)
		go func(name string, allData [][]byte) {
			err := readChunk(name, allData)
			if err != nil {
				log.Fatalf("Could not read chunk %v, error: %v\n", name, err)
			}
			wg.Done()
		}(name, allData)
	}
	wg.Wait()
	concatData := make([]byte, 0)
	for _, s := range allData {
		concatData = append(concatData, s...)
	}
	return concatData, nil
}

func readChunk(chunkPath string, allData [][]byte) error {
	filePath := filepath.Join("./.data", chunkPath)
	log.Printf("Opening Chunk %v\n", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	i, idxBytes := 0, make([]byte, 0)
	for {
		if data[i] == byte('\n') {
			break
		}
		idxBytes = append(idxBytes, data[i])
		i++
	}
	data = data[i+1:]
	idx, _ := strconv.Atoi(string(idxBytes))
	allData[idx] = data
	return nil
}
