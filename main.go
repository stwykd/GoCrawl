package main

import (
	"fmt"
	"os"
	"sync"
	"webCrawler/webcrawler"
)

// $ go run main.go seed
// if `seed` is not specified, it will be set to "https://monzo.com"
func main() {
	var seed = "https://monzo.com"
	if len(os.Args) == 2 {
		seed = os.Args[1]
	}

	wc, err := webcrawler.NewWebCrawler(seed); if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	numCrawlers := 50 // need at least two go routines for WebCrawl() as it reads and writes to the same channel
	wg.Add(numCrawlers)
	for i := 0; i < numCrawlers; i++ {
		go wc.WebCrawl(&wg)
	}
	wg.Wait()
	fmt.Println(wc.GetCrawledURLs())
}
