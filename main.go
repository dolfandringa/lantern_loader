package main

import (
	"fmt"
	"log"
	"os"
	"time"
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
	fmt.Println(urls)
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

}
