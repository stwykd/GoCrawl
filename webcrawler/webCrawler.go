package webcrawler

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var crawled = make(map[string]bool) // custom set using map for quick look-up
var toCrawl = make(chan string, 1000)
var mutex = sync.Mutex{}


type WebCrawler struct {
	seed string
}

func NewWebCrawler(seed string) (*WebCrawler, error) {
	if !isUrlValid(seed) {
		return &WebCrawler{}, errors.New("invalid seed url")
	}
	wc := WebCrawler{seed:seed}
	toCrawl <- seed
	return &wc, nil
}

func isUrlValid(url string) bool {
	 b, err := regexp.MatchString(`^http[s]?:[/][/][www.]?[\S]+[.][\S]+[.uk]?$`, url); if err != nil {
		panic(err)
	 }
	return b
}

// WebCrawl fetches all hyperlinks within the seed url, it will crawl all pages within
// the domain of the url without following external links Given a URL, and print a simple
// site map, showing the links between pages.
func (wc *WebCrawler) WebCrawl(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// read next url from to be crawled
		case url, ok := <-toCrawl:
			if ok {
				// fetch http from the next url in `toCrawl`
				html, err := fetchHttpFromUrl(url)
				if err != nil {
					panic(err)
				}

				// get all relative urls within the html
				urls := extractRelativeURLs(html)

				// add extracted urls to `toCrawl`
				for _, u := range urls {
					u := wc.seed+string(u)

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
				log.Println("the channel `toCrawl` has been closed")
				return
			}
			case <- time.After(4*time.Second):
				log.Println("Timeout reached, closing go routine")
				return
		}
	}
}

func fetchHttpFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), err
	}

	err = resp.Body.Close()
	if err != nil {
		return make([]byte, 0), err
	}
	return body, nil
}

func extractRelativeURLs(html []byte) []string {
	r := regexp.MustCompile(`href="(/[^".]+)"`)

	urls := make([]string, 0)
	for _, s := range r.FindAllStringSubmatch(string(html), -1) {
		urls = append(urls, s[1])
	}

	return urls
}

//func extractHttpHyperlinks(html []byte) [][]byte {
//	//r := regexp.MustCompile("<\\s*a\\s+[\\s\\S]*href\\s*=\\s*\"\\s*(http[^\"]*)\"[\\s\\S]*>[\\s\\S]*<\\s*/a\\s*>")
//	r := regexp.MustCompile("href\\s*=\\s*\"\\s*(http[^\"]*)\\s*\"")
//
//	// FindAllSubmatch() returns both the whole-pattern matches and the submatches within those matches.
//	// For example, using regex `r`, r.FindAllSubmatch("...<a hre="https://google.com">google.com</a>...") will
//	// return [[["...<a hre="https://google.com">google.com</a>...", "https://google.com"]]]
//	submatches := r.FindAllSubmatch(html, -1)
//
//	// remove whole-pattern matches from `submatches`
//	links := make([][]byte, 0)
//	for _, s := range submatches {
//		links = append(links, s[1])
//	}
//
//	return links
//}

func (wc *WebCrawler) GetCrawledURLs() []string {
	urls := make([]string, 0, len(crawled))
	for k := range crawled {
		urls = append(urls, k)
	}
	return urls
}