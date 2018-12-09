package sitemap

import (
	"strings"
)

type webPage struct {
	url      string
	subPages map[string]webPage
}

//func sortedMapKeys(m map[string]webPage) []string{
//	ks := make([]string, len(m))
//	for k := range m {
//		ks = append(ks, k)
//	}
//	sort.Strings(ks)
//	return ks
//}

func (w *webPage) string(tabs int) string {
	sb := strings.Builder{}
	sb.WriteString(strings.Repeat("\t", tabs)+w.url+"\n")
	for _, wp := range w.subPages {
			sb.WriteString(wp.string(tabs+1))
	}
	return sb.String()
}

func (w *webPage) String() string {
	return "\n\n\n"+w.string(1)
}