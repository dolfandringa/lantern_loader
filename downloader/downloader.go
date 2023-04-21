package downloader

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/schollz/progressbar/v3"

	"lantern_loader/utils"
)

func Downloader(urls []string) {
	//TODO add chunkSize argument
	parent_ctx := context.Background()
	ctx, cancel := context.WithCancel(parent_ctx)
	wg_size := sync.WaitGroup{}

	sizeChan := make(chan FileInfo)
	for _, url := range urls {
		wg_size.Add(1)
		go func(sizeChan chan<- FileInfo, url string, ctx context.Context) {
			defer wg_size.Done()
			fileInfo, err := GetFileInfo(url)
			select {
			case <-ctx.Done():
				return
			default:
				if err == nil {
					sizeChan <- fileInfo
				}
			}
		}(sizeChan, url, ctx)
	}
	fileInfo := <-sizeChan
	fmt.Println("File size: ", fileInfo.size)
	cancel()
	wg_size.Wait()
	const chunkSize = 1024 * 1024 * 10
	numChunks := int(math.Ceil(float64(fileInfo.size) / float64(chunkSize)))
	chunkChan := make(chan Job, numChunks)
	errorChan := make(chan Job, numChunks)
	writerChan := make(chan Chunk, numChunks)
	writerErrorChan := make(chan Chunk, numChunks)
	ctx, cancel = context.WithCancel(parent_ctx)
	wg := sync.WaitGroup{}
	var i int64
	for i = 0; i < fileInfo.size; i += chunkSize {
		chunkChan <- Job{Start: i, Stop: int64(utils.Min(i+chunkSize-1, fileInfo.size))}
	}
	fmt.Println("Dowloading ", len(chunkChan), "chunks")
	wg.Add(1)
	go func() {
		defer wg.Done()
		FileWriter(fileInfo.filename, writerChan, writerErrorChan, ctx)
	}()

	bar := progressbar.DefaultBytes(fileInfo.size)
	bar.Clear()
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			DownloadWorker(url, 10, chunkChan, errorChan, writerChan, ctx, bar)
		}(url)
	}
	for (len(chunkChan) > 0 || len(errorChan) > 0 || len(writerChan) > 0) && ctx.Err() == nil {
		select {
		case errorJob := <-errorChan:
			if errorJob.Retries < 5 {
				errorJob.Retries += 1
				chunkChan <- errorJob
			} else {
				cancel()
				fmt.Println("Failed downloading chunk ", errorJob, " 5 times.")
			}
		case writeError := <-writerErrorChan:
			fmt.Println("Error writing chunk", writeError.Start)
			cancel()
		}

	}
	//time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
}
