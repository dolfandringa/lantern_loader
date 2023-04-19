package main

import (
	"context"
	"log"
	"math"
	"os"

	"lantern_loader/downloader"
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
	var fileSize int = 10e9
	const chunkSize = 1024 * 1024
	chunkChan := make(chan downloader.Job, len(urls))
	errorChan := make(chan downloader.Job, len(urls))
	writerChan := make(chan []byte, len(urls))
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	for i := 0; i < fileSize; i += chunkSize {
		chunkChan <- downloader.Job{i, int(math.Min(float64(i+chunkSize-1), float64(fileSize)))}
	}
	// TODO: Fetch file size and name
	// Start writer

	for _, url := range urls {
		go downloader.DownloadWorker(url, 10, chunkChan, errorChan, writerChan, ctx)
	}

}
