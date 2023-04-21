package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func GetFileInfo(uri string) (FileInfo, error) {
	//TODO, should actually check Accept-Ranges header here, to see what ranges are accepted and return that in FileInfo
	req, err := http.NewRequest("HEAD", uri, nil)
	if err != nil {
		return FileInfo{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FileInfo{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == 405 {
		req, err = http.NewRequest("GET", uri, nil)
		if err != nil {
			return FileInfo{}, err
		}
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return FileInfo{}, err
		}
	}
	if res.StatusCode != 200 {
		return FileInfo{}, fmt.Errorf("We received status code %v", res.StatusCode)
	}
	if len(res.Header["Content-Length"]) == 0 {
		return FileInfo{}, errors.New("No Content-Length header present.")
	}
	sizeStr := res.Header["Content-Length"][0]
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return FileInfo{}, err
	}
	u, _ := url.Parse(uri)
	parts := strings.Split(u.Path, "/")
	filename := parts[len(parts)-1]
	return FileInfo{size: size, filename: filename}, nil
}
