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
				// fetch http from the next url in `toCrawl`
				html, err := utils.FetchHttpFromUrl(url)
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

func (wc *webCrawler) SpawnWebCrawlers(numCrawlers int) error {
	if numCrawlers < 2 {
		return errors.New("need at least two go routines for WebCrawl() as it reads and writes to the same channel")
	}
	var wg sync.WaitGroup
	wg.Add(numCrawlers)
	for i := 0; i < numCrawlers; i++ {
		go wc.WebCrawl(&wg)
	}
	wg.Wait()

	return nil
}

func (wc *webCrawler) GetCrawledURLs() []string {
	urls := make([]string, 0, len(crawled))
	for k := range crawled {
		urls = append(urls, k)
	}
	return urls
}

func (wc *webCrawler) PrintSitemap() {
	s, err := sitemap.NewSitemap(wc.seed, wc.GetCrawledURLs())
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}

//func (wc *webCrawler) PrintSitemap() {
//	m := map[int][]string{}
//	for _, url := range wc.GetCrawledURLs() {
//		n := strings.Count(url, "/")
//		m[n] = append(m[n], url)
//	}
//
//	var keys []int
//	for k := range m {
//		keys = append(keys, k)
//	}
//	sort.Ints(keys)
//
//	for _, k := range keys {
//		for _, url := range m[k] {
//			fmt.Println(strings.Repeat("\t", k) + url[len(wc.seed):])
//		}
//	}
//}
