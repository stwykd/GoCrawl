package utils

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

func Test_IsUrlValid(t *testing.T) {
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
			if got := IsUrlValid(tt.args.url); got != tt.want {
				t.Errorf("isUrlValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FetchHtmlFromUrl(t *testing.T) {
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
			got, err := FetchHtmlFromUrl(tt.args.url)
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

func Test_ExtractRelativeURLs(t *testing.T) {
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
			if got := ExtractRelativeURLs(tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}