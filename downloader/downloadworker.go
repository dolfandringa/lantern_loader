package downloader

import (
	"context"
	"fmt"
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
	fmt.Println("Starting worker for url", url)
	for {
		select {
		case job = <-jobChan:
		case <-ctx.Done():
			fmt.Println("Cancelled", url)
			return
		}
		fmt.Println("Fetching job", job)
		chunk, err := DownloadChunk(job.Start, job.Stop, url)
		if err != nil {
			fmt.Println("Error while downlinking", err)
			errorChan <- job
			continue
		}
		writerChan <- chunk
	}
}
