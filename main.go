package main

import (
	"os"
	"webCrawler/webcrawler"
)

// $ go run main.go seed
// if `seed` is not specified, `seed` will be set to "https://monzo.com"
// for example, to crawl facebook.com, run:
// $ go run main.go facebook.com
func main() {
	var seed = "https://monzo.com"
	if len(os.Args) == 2 {
		seed = os.Args[1]
	}

	wc, err := webcrawler.NewWebCrawler(seed)
	if err != nil {
		panic(err)
	}

	err = wc.SpawnWebCrawlers(1000)
	if err != nil {
		panic(err)
	}

	wc.PrintSitemap()
}
