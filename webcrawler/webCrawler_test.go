package webcrawler

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

var exampleComHtml = `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 50px;
        background-color: #fff;
        border-radius: 1em;
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        body {
            background-color: #fff;
        }
        div {
            width: auto;
            margin: 0 auto;
            border-radius: 0;
            padding: 1em;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is established to be used for illustrative examples in documents. You may use this
    domain in examples without prior coordination or asking for permission.</p>
    <p><a href="http://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>
`

func TestNewWebCrawler(t *testing.T) {
	type args struct {
		seed string
	}
	tests := []struct {
		name    string
		args    args
		want    *WebCrawler
		wantErr bool
	}{
		{
			"construct webcrawler with valid seed",
			args{"https://example.com"},
			&WebCrawler{"https://example.com"},
			false,
		},
		{
			"construct webcrawler with seed with invalid seed",
			args{"xyz"},
			&WebCrawler{},
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

func Test_isUrlValid(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"valid .com url",
			args{"https://example.com"},
			true,
		},
		{
			"valid .it url",
			args{"https://italia.it"},
			true,
		},
		{
			"valid .uk url",
			args{"https://example.ac.uk"},
			true,
		},
		{
			"url with no protocol",
			args{"example.com"},
			false,
		},
		{
			"url with incomplete protocol",
			args{"http:/example.com"},
			false,
		},
		{
			"invalid url",
			args{"xyz"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUrlValid(tt.args.url); got != tt.want {
				t.Errorf("isUrlValid() = %v, want %v", got, tt.want)
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

func Test_fetchHttpFromUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Fetch example.com",
			args{"https://example.com"},
			[]byte(exampleComHtml),
			false,
		},
		{"Fetch invalid url to get error",
			args{"xyz"},
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchHttpFromUrl(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("\ngot  %v\nwant %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}

func Test_extractRelativeURLs(t *testing.T) {
	type args struct {
		html []byte
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"simple html a tag",
			args{[]byte(`<a href="/blog">Blog</a>`)},
			[]string{"/blog"},
		},
		{
			"html a tag with subfolder",
			args{[]byte(`<a href="/static/images">Blog</a>`)},
			[]string{"/static/images"},
		},
		{
			"html a tag with protocol",
			args{[]byte(`<a href="https://www.example.com">Example</a>`)},
			[]string{},
		},
		{
			"multiple html a tags",
			args{[]byte(`<a href="/blog"></a><a href="/about"></a>`)},
			[]string{"/blog", "/about"},
		},
		{
			"html a tag with more attributes",
			args{[]byte(`<a name="blog" charset="UTF-8" href="/blog">Blog</a>`)},
			[]string{"/blog"},
		},
		{
			"html a tag to a local file",
			args{[]byte(`<a href ="/style.css`)},
			[]string{},
		},
		{
			"href with url only",
			args{[]byte(`href="/blog"`)},
			[]string{"/blog"},
		},
		{
			"href with url with no protocol",
			args{[]byte(`href="www.example.com"`)},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractRelativeURLs(tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}