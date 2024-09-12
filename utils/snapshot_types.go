package utils

type Snapshot struct {
	Filename string      `json:"filename"`
	Time     int64       `json:"time"`
	Size     int64       `json:"size"`
	Chunks   []ChunkInfo `json:"chunks"`
}

type ChunkInfo struct {
	FileName string `json:"name"`
	Order    int    `json:"order"`
}
