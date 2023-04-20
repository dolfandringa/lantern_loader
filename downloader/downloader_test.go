package downloader

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestDownloadChunk(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	wantStart := 0
	wantStop := 10
	gock.New(wantUrl).
		MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
		Get("/").
		Reply(200).
		Body(bytes.NewBuffer(want))
	got, err := DownloadChunk(wantStart, wantStop, wantUrl)
	st.Expect(t, err, nil)
	st.Expect(t, want, got)
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
	wantStart := 0
	wantStop := 10
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gock.New(wantUrl).
				MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
				Get("/").
				Reply(tc.status).
				Body(bytes.NewBuffer(want))
			_, err := DownloadChunk(wantStart, wantStop, wantUrl)
			st.Reject(t, err, nil)
		})
	}
}

func TestDownloadChunkWrongUrl(t *testing.T) {
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	defer gock.Off()
	wantUrl := "https://fake.url/file"
	wantStart := 0
	wantStop := 10
	gock.New(wantUrl).
		MatchHeader("Range", fmt.Sprintf("bytes=%d-%d", wantStart, wantStop)).
		Get("/bla").
		Reply(200).
		Body(bytes.NewBuffer(want))
	_, err := DownloadChunk(wantStart, wantStop, wantUrl)
	st.Reject(t, err, nil)
}
