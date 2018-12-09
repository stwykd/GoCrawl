// Package utils provides helper methods used throughout the project for common web related tasks
package utils

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

// IsUrlValid checks if a url is valid. Valid urls use the http protocol, may have a domain starting
// with www, have a domain name comprised of letters and numbers, and ends with a top-level domain
func IsUrlValid(url string) bool {
	b, err := regexp.MatchString(`^http[s]?:[/][/][www.]?[\S]+[.][\S]+[.uk]?$`, url)
	if err != nil {
		panic(err)
	}
	return b
}

// FetchHtmlFromUrl fetches the html code from a url.
// An error is returned if the http.Get() function fails
func FetchHtmlFromUrl(url string) ([]byte, error) {
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

// ExtractRelativeURLs extracts all relative URLs from html code
func ExtractRelativeURLs(html []byte) []string {
	r := regexp.MustCompile(`href="(/[^".]+)"`)
	submatches := r.FindAllStringSubmatch(string(html), -1)

	n := len(submatches)
	urls := make([]string, n)
	for i := 0; i < n; i++ {
		urls[i] = submatches[i][1]
	}

	return urls
}

//func ExtractHttpHyperlinks(html []byte) [][]byte {
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

// RemoveHttpFromURL removes the http:// or https:// from the beginning of an url
func RemoveHttpFromURL(url string) string {
	return regexp.MustCompile(`^http[s]?://([\S]+)$`).ReplaceAllString(url, "${1}")
}
