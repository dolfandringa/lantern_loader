package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func DownloadChunk(start, end int64, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && res.StatusCode != 206 {
		return nil, fmt.Errorf("We received status code %v instead of 200", res.StatusCode)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil

}

func GetSize(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode == 405 {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, err
		}
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
	}
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("We received status code %v", res.StatusCode)
	}
	if len(res.Header["Content-Length"]) == 0 {
		return 0, errors.New("No Content-Length header present.")
	}
	sizeStr := res.Header["Content-Length"][0]
	size, err := strconv.ParseInt(sizeStr, 10, 0)
	if err != nil {
		return 0, err
	}
	return int64(size), nil
}
