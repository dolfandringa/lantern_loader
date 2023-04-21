package downloader

import (
	"context"
	"fmt"

	"github.com/schollz/progressbar/v3"
)

func DownloadWorker(
	url string,
	timeout int,
	jobChan <-chan Job,
	errorChan chan<- Job,
	writerChan chan<- Chunk,
	ctx context.Context,
	bar *progressbar.ProgressBar,
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
		chunk, err := DownloadChunk(job.Start, job.Stop, url)
		if err != nil {
			fmt.Println("Error while downlinking", err)
			errorChan <- job
			break
		}
		_ = bar.Add64(int64(job.Stop - job.Start))
		writerChan <- Chunk{data: chunk, Start: job.Start, Stop: job.Stop}
	}
}
