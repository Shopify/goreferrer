package referrer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRelativeUrl(t *testing.T) {
	url := `buh`
	r, err := Parse(url)

	assert.NoError(t, err)
	indirect := r.(*Indirect)
	assert.Equal(t, indirect.Url, url)
}

func TestNotSearchDirectOrSocial(t *testing.T) {
	url := "http://unicorns.ca/"
	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, url, r.(*Indirect).Url)
}

func TestSearchSimple(t *testing.T) {
	r, err := Parse("http://ca.search.yahoo.com/search?p=hello")
	assert.NoError(t, err)
	switch r := r.(type) {
	case *Search:
		assert.Equal(t, r.Label, "Yahoo")
		assert.Equal(t, r.Query, "hello")
	default:
		assert.Fail(t, "Wrong referrer result!")
	}
}

func TestSearchNonAscii(t *testing.T) {
	url := "http://ca.search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1"
	assert.True(t, strings.Contains(url, "\u00F8"))

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "\u00F8"))
	assert.Equal(t, engine.Query, "vinduespudsning myshopify rengøring mkobetic")
}

func TestSearchWithExplicitPlus(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+JOKAPOLAR+"11+%2B+11"+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "11 + 11"))
	assert.Equal(t, engine.Query, `vinduespudsning JOKAPOLAR "11 + 11" mkobetic`)
}

func TestSearchWithNonAscii(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "rengøring"))
	assert.Equal(t, engine.Query, `vinduespudsning myshopify rengøring mkobetic`)
}

func TestSearchWithCyrillics(t *testing.T) {
	url := `http://www.yandex.com/yandsearch?text=%D0%B1%D0%BE%D1%82%D0%B8%D0%BD%D0%BA%D0%B8%20packer-shoes&lr=87&msid=22868.18811.1382712652.60127&noreask=1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yandex")
	assert.True(t, strings.Contains(engine.Query, "ботинки"))
	assert.Equal(t, engine.Query, `ботинки packer-shoes`)
}

func TestDirectSimple(t *testing.T) {
	url := "http://example.com"

	r, err := ParseWithDirect(url, "example.com", "sample.com")
	assert.NoError(t, err)

	direct := r.(*Direct)
	assert.NotNil(t, direct)
	assert.Equal(t, direct.Url, url)
	assert.Equal(t, direct.Domain, "example.com")
}

func TestSocialSimple(t *testing.T) {
	url := "https://twitter.com/snormore/status/391149968360103936"

	r, err := Parse(url)
	assert.NoError(t, err)

	social := r.(*Social)
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Twitter")
}

func ExampleParse() {
	urls := []string{
		"http://ca.search.yahoo.com/search?p=hello",
		"https://twitter.com/jdoe/status/391149968360103936",
		"http://yoursite.com/links",
	}

	for _, url := range urls {
		r, err := Parse(url)
		if err != nil {
			panic(err)
		}
		switch r := r.(type) {
		case *Search:
			fmt.Printf("Search %s: %s\n", r.Label, r.Query)
		case *Social:
			fmt.Printf("Social %s\n", r.Label)
		case *Indirect:
			fmt.Printf("Indirect: %s\n", r.Url)
		}
	}
	// Output:
	// Search Yahoo: hello
	// Social Twitter
	// Indirect: http://yoursite.com/links
}
