package sitemap

import (
	"errors"
	"fmt"
	"strings"
	"webCrawler/utils"
)

type sitemap struct {
	domain webPage
}

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
			return nil, errors.New("url " + url + "is not in domain " + domain)
		}

		curr := s.domain
		for _, page := range splitUrl[1:] {
			if _, ok := curr.subPages[page]; !ok {
				curr.subPages[page] = webPage{page, make(map[string]webPage)}
			}
			curr = curr.subPages[page]
		}
	}

	return &s, nil
}

func (s *sitemap) String() string {
	b := strings.Builder{}
	b.WriteString(s.domain.url)
	return fmt.Sprint(b.String())
}
