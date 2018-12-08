package sitemap

import "fmt"

type webPage struct {
	url      string
	subPages map[string]webPage
}

func (w *webPage) String() string {
	return fmt.Sprintf("%v -> %v", w.url, w.subPages)
}
