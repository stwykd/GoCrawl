package webcrawler

import (
	"reflect"
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
			args{"https://example.it"},
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

func TestWebCrawler_PrintMap(t *testing.T) {
	tests := []struct {
		name string
		wc   *WebCrawler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wc.PrintMap()
		})
	}
}
