package main

import (
	"os"
	"webCrawler/webcrawler"
)

// $ go run main.go seed
// If `seed` is not specified, `seed` will be set to "https://monzo.com"
// Make sure `seed` has an http protocol. monzo.com will not work
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
