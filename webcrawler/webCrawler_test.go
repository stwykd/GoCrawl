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

func TestWebCrawler_WebCrawl(t *testing.T) {
	tests := []struct {
		name string
		w    *WebCrawler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.WebCrawl()
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
		{	"Fetch example.com",
			args{"https://example.com"},
			[]byte(exampleComHtml),
			false,
		},
		{	"Fetch invalid url to get error",
			args{"https://www.monzo.com/"},
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

func Test_extractHttpHyperlinks(t *testing.T) {
	type args struct {
		html []byte
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		{
			"Simple html a tag",
			args{[]byte(`<a href="https://www.example.com">Example</a>`)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"Multiple html a tags",
			args{[]byte(`<a href="https://www.example.com"">Example</a>` +
				`<a href="https://www.monzo.com">Monzo</a>`)},
			[][]byte{[]byte("https://www.example.com"), []byte("https://www.monzo.com")},
		},
		{
			"html a tag with spaces",
			args{[]byte(`<   a href  = "https://www.example.com">Example</a>`)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"html a tag with spaces in href url",
			args{[]byte(`<a href="   https://www.example.com">Example</a> `)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"html a tag with more attributes",
			args{[]byte(`<a name="example" charset="UTF-8" href="https://www.example.com">Example</a>`)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"html a tag with local url",
			args{[]byte(`<a href ="/home.html`)},
			[][]byte{},
		},
		{
			"href with url only",
			args{[]byte(`href="https://www.example.com"`)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"href with url with http protocol",
			args{[]byte(`href="http://www.example.com"`)},
			[][]byte{[]byte("http://www.example.com")},
		},
		{
			"href with url with https protocol",
			args{[]byte(`href="https://www.example.com"`)},
			[][]byte{[]byte("https://www.example.com")},
		},
		{
			"href with url with file protocol",
			args{[]byte(`href="file://www.example.com"`)},
			[][]byte{},
		},
		{
			"href with url with no protocol",
			args{[]byte(`href="www.example.com"`)},
			[][]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractHttpHyperlinks(tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}
