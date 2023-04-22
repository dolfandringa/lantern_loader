package downloader

import (
	"context"
	"fmt"
	"io"
	"os"
)

type fileSystem interface {
	Create(name string) (file, error)
}

type file interface {
	io.WriterAt
	io.Closer
}

// Mocker friendly os filesystem implementation.
type osFS struct{}

func (*osFS) Create(filename string) (file, error) { return os.Create(filename) }

var dfFS fileSystem = &osFS{}

func FileWriter(
	filename string,
	chunkChan <-chan Chunk,
	errorChan chan<- Chunk,
	completeChan chan<- Chunk,
	ctx context.Context,
) {
	//TODO handle existing file.
	file, err := dfFS.Create(filename)
	if err != nil {
		fmt.Println("Error opening file", filename, err)
		return
	}
	defer file.Close()
	for ctx.Err() == nil {
		select {
		case chunk := <-chunkChan:
			_, err := file.WriteAt(chunk.data, chunk.Start)
			if err != nil {
				fmt.Println(err)
				errorChan <- chunk
			}
			completeChan <- chunk
		case <-ctx.Done():
			fmt.Println("Cancelled writer")
			return
		}
	}
}
