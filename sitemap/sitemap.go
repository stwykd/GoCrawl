package sitemap

import (
	"errors"
	"strings"
	"webCrawler/utils"
)

type sitemap struct {
	domain webPage
}

// NewSitemap constructs a sitemap for domain `domain` using URLs `urls`.
// Returns an error if `domain` is not a valid URL or if any URL within `urls` is not in domain `domain`
func NewSitemap(domain string, urls []string) (*sitemap, error) {
	if !utils.IsUrlValid(domain) {
		return &sitemap{}, errors.New("invalid domain url: " + domain)
	}
	domain = utils.RemoveHttpFromURL(domain)
	s := sitemap{domain: webPage{domain, map[string]webPage{}}}

	for _, url := range urls {
		url := utils.RemoveHttpFromURL(url)
		splitUrl := strings.Split(url, "/")
		if splitUrl[0] != domain {
			return &sitemap{}, errors.New("url " + url + "is not in domain " + domain)
		}

		curr := s.domain
		splitUrl = splitUrl[1 : len(splitUrl)-1] // remove domain name and trailing whitespace
		for _, page := range splitUrl {
			if _, ok := curr.subPages[page]; !ok {
				curr.subPages[page] = webPage{page, make(map[string]webPage)}
			}
			curr = curr.subPages[page]
		}
	}

	return &s, nil
}

func (s *sitemap) String() string {
	return s.domain.String()
}
