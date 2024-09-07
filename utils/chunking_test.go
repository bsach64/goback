package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSize(t *testing.T) {
	testDir := "../test_files"
	err := filepath.Walk(testDir, walkTestFiles(t))
	if err != nil {
		t.Fatalf("Error walking through test files: %v", err)
	}
}

func walkTestFiles(t *testing.T) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				testFile(t, path)
			})
		}

		return nil
	}
}

func testFile(t *testing.T, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("File not found %v", err)
	}
	defer file.Close()

	stat, statErr := file.Stat()
	if statErr != nil {
		t.Fatalf("Cannot get stat for file")
	}

	chunkedFile, err := ChunkFile(filePath)
	if err != nil {
		t.Fatalf("Chunking failed: %v", err)
	}

	if chunkedFile.Meta.Size != stat.Size() {
		t.Fatalf("Size not same. Actual Size: %d, Chunked total size: %d", stat.Size(), chunkedFile.Meta.Size)
	}
}
