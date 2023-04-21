package downloader

type Chunk struct {
	Start int64
	Stop  int64
	data  []byte
}

type Job struct {
	Start   int64
	Stop    int64
	Retries int
}
