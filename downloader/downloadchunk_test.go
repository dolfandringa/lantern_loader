package downloader

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestDownloadChunk(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	var wantStart int64 = 0
	var wantStop int64 = 10
	gock.New(wantUrl).
		MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
		Get("/").
		Reply(200).
		Body(bytes.NewBuffer(want))
	got, err := DownloadChunk(wantStart, wantStop, wantUrl)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}

func TestDownloadChunkWith206(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	var wantStart int64 = 0
	var wantStop int64 = 10
	gock.New(wantUrl).
		MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
		Get("/").
		Reply(206).
		Body(bytes.NewBuffer(want))
	got, err := DownloadChunk(wantStart, wantStop, wantUrl)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}

func TestDownloadChunkHTTPError(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	cases := []struct {
		name   string
		status int
	}{
		{
			name:   "NotFound",
			status: 404,
		},
		{name: "InternalServerError",
			status: 500,
		}, {
			name:   "Unauthorized",
			status: 401,
		}, {
			name:   "Forbidden",
			status: 403,
		},
	}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	var wantStart int64 = 0
	var wantStop int64 = 10
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gock.New(wantUrl).
				MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
				Get("/").
				Reply(tc.status).
				Body(bytes.NewBuffer(want))
			_, err := DownloadChunk(wantStart, wantStop, wantUrl)
			assert.NotEqual(t, err, nil)
		})
	}
}

func TestDownloadChunkWrongUrl(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	var wantStart int64 = 0
	var wantStop int64 = 10
	gock.New(wantUrl).
		MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
		Get("/bla").
		Reply(200).
		Body(bytes.NewBuffer(want))
	_, err := DownloadChunk(wantStart, wantStop, wantUrl)
	assert.NotEqual(t, err, nil)
}

func TestGetFileInfoMissingContentLength(t *testing.T) {
	wantUrl := "https://fake.url/file"
	defer gock.Off()
	gock.New(wantUrl).
		Head("/").
		Reply(200)
	_, err := GetFileInfo(wantUrl)
	assert.NotEqual(t, err, nil)
}

func TestGetFileInfo(t *testing.T) {
	wantUrl := "https://fake.url/file.txt#test?q=1&r=2"
	want := FileInfo{size: 50000, filename: "file.txt"}
	defer gock.Off()
	gock.New(wantUrl).
		Head("/").
		Reply(200).
		AddHeader("Content-Length", fmt.Sprint(want.size))
	got, err := GetFileInfo(wantUrl)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}

func TestGetFileInfoCatch405(t *testing.T) {
	wantUrl := "https://fake.url/file.txt"
	want := FileInfo{size: 50000, filename: "file.txt"}
	defer gock.Off()
	gock.New(wantUrl).
		Get("/").
		Reply(200).
		AddHeader("Content-Length", fmt.Sprint(want.size))
	gock.New(wantUrl).
		Head("/").
		Reply(405)
	got, err := GetFileInfo(wantUrl)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}
