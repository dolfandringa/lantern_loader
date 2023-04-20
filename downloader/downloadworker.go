package downloader

import (
	"context"
	"time"
)

func DownloadWorker(
	url string,
	timeout int,
	jobChan <-chan Job,
	errorChan chan<- Job,
	writerChan chan<- []byte,
	ctx context.Context,
) {
	var job Job
	for {
		select {
		case job = <-jobChan:
		case <-time.After(time.Duration(timeout) * time.Second):
			return
		case <-ctx.Done():
			return
		}
		chunk, err := DownloadChunk(job.Start, job.Stop, url)
		if err != nil {
			errorChan <- job
		}
		writerChan <- chunk
	}
}
