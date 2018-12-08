package main

import (
	"os"
	"strconv"
	"sync"
	"webCrawler/webcrawler"
)

// $ go run main.go seed numCrawlers
// if `seed` and `numCrawlers` are both not specified, `seed` will be set to "https://monzo.com"
// and `numCrawlers` will be se to 3
func main() {
	var seed, numCrawlers = "https://monzo.com", 3
	if len(os.Args) == 3 {
		seed = os.Args[1]
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		numCrawlers = n
	}

	wc := webcrawler.NewWebCrawler(seed)
	var wg sync.WaitGroup
	wg.Add(numCrawlers)
	for i := 0; i < numCrawlers; i++ {
		go wc.WebCrawl(&wg)
	}

	wc.PrintMap()
}
