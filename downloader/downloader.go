package downloader

import (
	"context"
	"fmt"
	"math"
	"sync"

	"lantern_loader/utils"
)

func Downloader(urls []string) {
	parent_ctx := context.Background()
	ctx, cancel := context.WithCancel(parent_ctx)
	wg_size := sync.WaitGroup{}

	sizeChan := make(chan int64)
	for _, url := range urls {
		wg_size.Add(1)
		go func(sizeChan chan<- int64, url string, ctx context.Context) {
			defer wg_size.Done()
			size, err := GetSize(url)
			select {
			case <-ctx.Done():
				return
			default:
				if err == nil {
					sizeChan <- size
				}
			}
		}(sizeChan, url, ctx)
	}
	fileSize := <-sizeChan
	fmt.Println("File size: ", fileSize)
	cancel()
	wg_size.Wait()
	const chunkSize = 1024
	numChunks := int(math.Ceil(float64(fileSize) / float64(chunkSize)))
	chunkChan := make(chan Job, numChunks)
	errorChan := make(chan Job, numChunks)
	writerChan := make(chan []byte, numChunks)
	ctx, cancel = context.WithCancel(parent_ctx)
	var i int64
	for i = 0; i < fileSize; i += chunkSize {
		chunkChan <- Job{Start: i, Stop: int64(utils.Min(i+chunkSize-1, fileSize))}
	}
	fmt.Println("Dowloading ", len(chunkChan), "chunks")
	// TODO: Fetch file name
	// Start writer

	for _, url := range urls {
		go DownloadWorker(url, 10, chunkChan, errorChan, writerChan, ctx)
	}
	for (len(chunkChan) > 0 || len(errorChan) > 0) && ctx.Err() == nil {
		var errorJob Job
		for len(errorChan) > 0 {
			errorJob = <-errorChan
			if errorJob.Retries < 5 {
				errorJob.Retries += 1
				chunkChan <- errorJob
			} else {
				cancel()
				fmt.Println("Failed downloading chunk ", errorJob, " 5 times.")
			}
		}

	}
	cancel()
	fmt.Println(
		"ctx.Err():",
		ctx.Err(),
		"remaining chunks: ",
		len(chunkChan),
		"remaining errors:",
		len(errorChan),
	)
}
