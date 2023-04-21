package main

import (
	"fmt"
	"log"
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
	downloader.Downloader(urls)
	fmt.Println("Done downloading")
}
