package webcrawler

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestNewWebCrawler(t *testing.T) {
	type args struct {
		seed string
	}
	tests := []struct {
		name    string
		args    args
		want    *webCrawler
		wantErr bool
	}{
		{
			"construct webcrawler with valid seed",
			args{"https://example.com"},
			&webCrawler{"https://example.com"},
			false,
		},
		{
			"construct webcrawler with seed with invalid seed",
			args{"xyz"},
			&webCrawler{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebCrawler(tt.args.seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebCrawler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebCrawler_WebCrawl(t *testing.T) {
	want := []string{"https://google.co.uk/accounts/answer/39612", "https://google.co.uk/?hl=en-GB", "https://google.co.uk/accounts", "https://google.co.uk/chrome/answer/95647", "https://google.co.uk/accounts/?hl=en#topic=3382296", "https://google.co.uk/accounts/answer/58585?hl=en&amp;ref_topic=3382296", "https://google.co.uk/accounts/answer/3265955?hl=en&amp;ref_topic=3382296", "https://google.co.uk/preferences?hl=en", "https://google.co.uk/services/", "https://google.co.uk/history/optout?hl=en-GB", "https://google.co.uk/preferences?hl=en-GB", "https://google.co.uk/accounts/answer/27441?hl=en&amp;ref_topic=3382296", "https://google.co.uk/intl/en/policies/terms/", "https://google.co.uk/support/websearch?p=ws_cookies_notif&amp;hl=en-GB", "https://google.co.uk/accounts/answer/114129?hl=en&amp;ref_topic=3382296", "https://google.co.uk/accounts/troubleshooter/2402620?hl=en&amp;ref_topic=3382296", "https://google.co.uk/webhp", "https://google.co.uk/chrome/answer/95464", "https://google.co.uk/um/StartNow?hl=en&amp;sourceid=awo&amp;subid=ww-ww-di-g-aw-a-awhp_1!o2", "https://google.co.uk/accounts/answer/183723?hl=en&amp;ref_topic=3382296", "https://google.co.uk/advanced_search?hl=en-GB&amp;authuser=0", "https://google.co.uk/language_tools?hl=en-GB&amp;authuser=0", "https://google.co.uk/accounts/answer/3118687?hl=en&amp;ref_topic=3382296", "https://google.co.uk/intl/en/policies/privacy/", "https://google.co.uk/accounts/answer/32040?hl=en&amp;ref_topic=3382296", "https://google.co.uk/accounts/answer/32050", "https://google.co.uk/accounts/answer/6304920?hl=en&amp;ref_topic=3382296", "https://google.co.uk/intl/en/ads/", "https://google.co.uk/intl/en/policies/"}

	wc, err := NewWebCrawler("https://google.co.uk")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	numCrawlers := 20
	wg.Add(numCrawlers)
	for i := 0; i < numCrawlers; i++ {
		go wc.WebCrawl(&wg)
	}
	wg.Wait()

	got := wc.GetCrawledURLs()
	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nWebCrawl() crawled:\n%v\nbut was expected to crawl:\n%v\n", got, want)
	}

}
