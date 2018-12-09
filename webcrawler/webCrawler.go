// package webcrawler implements a concurrent web crawler that, starting from a seed url,
// discovers all pages within the same domain of the seed url
package webcrawler

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"webCrawler/sitemap"
	"webCrawler/utils"
)

var crawled = make(map[string]bool) // custom set using map for quick look-up
var toCrawl = make(chan string, 1000)
var mutex = sync.Mutex{}

type webCrawler struct {
	seed string
}

// NewWebCrawler constructs a new web crawler using a seed url.
// It returns an error if the seed url is not a valid url
func NewWebCrawler(seed string) (*webCrawler, error) {
	if !utils.IsUrlValid(seed) {
		return &webCrawler{}, errors.New("invalid seed url: " + seed)
	}
	wc := webCrawler{seed: seed}
	toCrawl <- seed
	return &wc, nil
}

// WebCrawl fetches all hyperlinks within the seed url, it will crawl all pages within
// the domain of the url without following external links Given a URL, and print a simple
// site map, showing the links between pages.
func (wc *webCrawler) WebCrawl(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// read the next url to be crawled from `toCrawl`
		case url, ok := <-toCrawl:
			if ok {
				// fetch html code from the next url in `toCrawl`
				html, err := utils.FetchHtmlFromUrl(url)
				if err != nil {
					panic(err)
				}

				// get all relative urls from the fetched html
				urls := utils.ExtractRelativeURLs(html)

				// add the extracted urls to `toCrawl`
				for _, u := range urls {
					u := wc.seed + string(u)

					mutex.Lock()
					if !crawled[u] {
						toCrawl <- u
						crawled[u] = true
						log.Printf("%v added to channel \n", u)
						log.Printf("crawled %v urls, %v left to crawl \n", len(crawled), len(toCrawl))
					}
					mutex.Unlock()
				}
			} else {
				log.Println("the channel `toCrawl` has been closed, killing go routine")
				return
			}
		case <-time.After(4 * time.Second):
			log.Println("Timeout reached, killing go routine")
			return
		}
	}
}

// SpawnWebCrawlers spawns numCrawlers web crawlers that concurrently discover all pages within the domain
// of the seed url
func (wc *webCrawler) SpawnWebCrawlers(numCrawlers int) error {
	if numCrawlers < 2 {
		return errors.New("need at least two goroutines for WebCrawl() as the function both reads from " +
			"and writes to the same channel")
	}
	var wg sync.WaitGroup
	wg.Add(numCrawlers)
	for i := 0; i < numCrawlers; i++ {
		go wc.WebCrawl(&wg)
	}
	wg.Wait()

	return nil
}

// GetCrawledURLs returns all URLs discovered by the web crawler
func (wc *webCrawler) GetCrawledURLs() []string {
	urls := make([]string, 0, len(crawled))
	for k := range crawled {
		urls = append(urls, k)
	}
	return urls
}

// PrintSitemap prints a sitemap of the crawled domain
func (wc *webCrawler) PrintSitemap() {
	s, err := sitemap.NewSitemap(wc.seed, wc.GetCrawledURLs())
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}
