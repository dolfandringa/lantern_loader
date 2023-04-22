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
	writerChan chan<- Chunk,
	ctx context.Context,
) {
	var job Job
	fmt.Println("Starting worker for url", url)
	for len(jobChan) > 0 {
		select {
		case job = <-jobChan:
		case <-ctx.Done():
			fmt.Println("Cancelled", url)
			return
		}
		chunk, err := DownloadChunk(job.Start, job.Stop, url)
		if err != nil {
			fmt.Println("Error while downlinking", err)
			errorChan <- job
			break
		}
		writerChan <- Chunk{data: chunk, Start: job.Start, Stop: job.Stop}
	}
}
