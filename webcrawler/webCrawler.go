package webcrawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var crawled = make(map[string]bool) // custom set using map for quick look-up
var toCrawl = make(chan string)

type WebCrawler struct {
	Seed string
}

func (*WebCrawler) WebCrawl() {
	select {
	// read next url from to be crawled
	case url, ok := <-toCrawl:
		if ok {
			// fetch http from the next url in `toCrawl`
			html, err := fetchHttpFromUrl(url)
			if err != nil {
				panic(err)
			}

			// get all hyperlinks within the html
			links := extractHttpHyperlinks(html)

			// add extracted url to `toCrawl`
			for _, l := range links {
				url := string(l)
				if !crawled[url] {
					toCrawl <- url
					crawled[url] = true
				}
			}
		} else {
			log.Print("the channel `toCrawl` is closed")
			return
		}
	default:
		log.Print("no more urls to crawl")
		return
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

	// cannot `defer resp.Body.Close()` because Close() returns an error that needs to checked
	err = resp.Body.Close()
	if err != nil {
		return make([]byte, 0), err
	}
	return body, nil
}

func extractHttpHyperlinks(html []byte) [][]byte {
	//r := regexp.MustCompile("<\\s*a\\s+[\\s\\S]*href\\s*=\\s*\"\\s*(http[^\"]*)\"[\\s\\S]*>[\\s\\S]*<\\s*/a\\s*>")
	r := regexp.MustCompile("href\\s*=\\s*\"\\s*(http[^\"]*)\\s*\"")

	// FindAllSubmatch() returns both the whole-pattern matches and the submatches within those matches.
	// For example, using regex `r`, r.FindAllSubmatch("...<a hre="https://google.com">google.com</a>...") will
	// return [[["...<a hre="https://google.com">google.com</a>...", "https://google.com"]]]
	submatches := r.FindAllSubmatch(html, -1)

	// remove whole-pattern matches from `submatches`
	links := make([][]byte, 0)
	for _, s := range submatches {
		links = append(links, s[1])
	}

	return links
}
