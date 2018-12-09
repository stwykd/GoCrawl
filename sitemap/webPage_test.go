package sitemap

import "testing"

func Test_webPage_String(t *testing.T) {
	tests := []struct {
		name string
		w    *webPage
		want string
	}{
		{
			"webpage with no subspages",
			&webPage{
				"example.com",
				map[string]webPage{},
			},
			"\texample.com\n",
		},
		{
			"webpage with one subpage",
			&webPage{
				"example.com",
				map[string]webPage{
					"about": {
						"about", map[string]webPage{},
					},
				},
			},
			"\texample.com\n\t\tabout\n",
		},
		{
			"webpage with subpage within subpage",
			&webPage{
				"example.com",
				map[string]webPage{
					"about": {
						"about", map[string]webPage{
							"our-story": {
								"our-story", map[string]webPage{},
							},
						},
					},
				},
			},
			"\texample.com\n\t\tabout\n\t\t\tour-story\n",
		},
		{
			"webpage with more subpages",
			&webPage{
				"example.com",
				map[string]webPage{
					"about": {
						"about", map[string]webPage{},
					},
					"legal": {
						"legal", map[string]webPage{},
					},
				},
			},
			"\texample.com\n\t\tabout\n\t\tlegal\n",
		},
		{
			"webpage with subpages within subpages",
			&webPage{
				"example.com",
				map[string]webPage{
					"about": {
						"about", map[string]webPage{
							"our-story": {
								"our-story", map[string]webPage{},
							},
						},
					},
					"legal": {
						"legal", map[string]webPage{
							"terms-and-conditions": {
								"terms-and-conditions", map[string]webPage{},
							},
						},
					},
				},
			},
			"\texample.com\n\t\tabout\n\t\t\tour-story\n\t\tlegal\n\t\t\tterms-and-conditions\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.String(); got != tt.want {
				t.Errorf("\nwebPage.String(): \n%v want: \n%v", got, tt.want)
			}
		})
	}
}
