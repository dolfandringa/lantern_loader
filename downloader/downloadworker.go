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
			return
		}
		fmt.Println("Fetching job", job)
		chunk, err := DownloadChunk(job.Start, job.Stop, url)
		fmt.Println("Finished fetching job", job)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Writing error", job)
			errorChan <- job
		}
		writerChan <- chunk
	}
}
