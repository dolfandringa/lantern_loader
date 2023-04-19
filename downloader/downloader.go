package downloader

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadChunk(start, end int, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil

}
