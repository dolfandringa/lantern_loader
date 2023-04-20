package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"lantern_loader/downloader"
	"lantern_loader/utils"
)

const HELP string = `
lantern_loader  [URL]...

This downloads the file present at all the different URL's. It is assumed the urls point to the same file.
The file name of the first request is used to store the file.
`

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		log.Fatal("No urls supplied")
	}
	parent_ctx := context.Background()
	ctx, cancel := context.WithCancel(parent_ctx)
	sizeChan := make(chan int)
	for _, url := range urls {
		go func(sizeChan chan<- int, url string, ctx context.Context) {
			size, err := downloader.GetSize(url)
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
	fmt.Println("Got size", fileSize)
	cancel()
	const chunkSize = 1024
	numChunks := int(math.Ceil(float64(fileSize) / float64(chunkSize)))
	chunkChan := make(chan downloader.Job, numChunks*2)
	errorChan := make(chan downloader.Job, numChunks*2)
	writerChan := make(chan []byte, numChunks)
	ctx, cancel = context.WithCancel(parent_ctx)
	for i := 0; i < fileSize; i += chunkSize {
		chunkChan <- downloader.Job{Start: i, Stop: int(utils.Min(i+chunkSize-1, fileSize))}
	}
	fmt.Println(len(chunkChan), "items put on the queue")
	// TODO: Fetch file name
	// Start writer

	for _, url := range urls {
		fmt.Println("Starting worker for url", url)
		go downloader.DownloadWorker(url, 10, chunkChan, errorChan, writerChan, ctx)
	}
	for len(chunkChan) > 0 {
		//fmt.Println(len(errorChan), "errors on the queue.")
		var errorJob downloader.Job
		for len(errorChan) > 0 {
			errorJob = <-errorChan
			fmt.Println("Got error for job", errorJob, "job queue length", len(chunkChan))
			if errorJob.Retries < 5 {
				errorJob.Retries += 1
				chunkChan <- errorJob
			} else {
				cancel()
				log.Fatalf("Failed downloading chunk %v 5 times.", errorJob)
			}
		}
		//fmt.Println(len(chunkChan), "items remaining on the queue")

	}
	fmt.Println("Done downloading")

}
