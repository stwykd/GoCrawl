package sitemap

import (
	"reflect"
	"testing"
)

func TestNewSitemap(t *testing.T) {
	type args struct {
		domain string
		urls   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *sitemap
		wantErr bool
	}{
		{
			"construct sitemap with valid domain and URLs",
			args{
				"https://example.com",
				[]string{"https://example.com/about", "https://example.com/legal"},
			},
			&sitemap{
				webPage{
					"example.com", map[string]webPage{},
				},
			},
			false,
		},
		{
			"construct sitemap with URLs outside of domain",
			args{
				"https://example.com",
				[]string{"https://google.com/about", "https://example.com/legal"},
			},
			&sitemap{},
			true,
		},
		{
			"construct sitemap with invalid domain",
			args{
				"https://example",
				[]string{"https://example.com/about", "https://example.com/legal"},
			},
			&sitemap{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSitemap(tt.args.domain, tt.args.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSitemap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSitemap() = %v, want %v", got, tt.want)
			}
		})
	}
}
