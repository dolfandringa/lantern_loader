package downloader

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockFS struct{ mock.Mock }
type mockFile struct{ mock.Mock }

func (o *mockFS) Create(name string) (file, error) {
	args := o.Called(name)
	return args.Get(0).(*mockFile), args.Error(1)
}

func (o *mockFile) WriteAt(p []byte, off int64) (int, error) {
	args := o.Called(p, off)
	return args.Int(0), args.Error(1)
}

func (o *mockFile) Close() error {
	args := o.Called()
	return args.Error(0)
}

func TestWriteChunk(t *testing.T) {
	mFS := new(mockFS)
	mF := new(mockFile)
	data := []byte("Hello World!")
	dfFS = mFS
	fname := "testfile.txt"
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mF.On("Close").Return(nil)
	mF.On("WriteAt", data, int64(10)).Return(len(data), nil)
	mFS.On("Create", fname).Return(mF, nil)
	chunkChan := make(chan Chunk, 2)
	errorChan := make(chan Chunk, 2)
	chunk := Chunk{Start: 10, Stop: int64(10 + len(data)), data: data}
	wg := sync.WaitGroup{}
	chunkChan <- chunk
	wg.Add(1)
	go func() {
		FileWriter(fname, chunkChan, errorChan, ctx)
		wg.Done()
	}()
	for len(chunkChan) > 0 {
	}
	cancel()
	wg.Wait()
	mF.AssertExpectations(t)
	mFS.AssertExpectations(t)
}
