package referrer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRelativeURL(t *testing.T) {
	url := `buh`
	r, err := Parse(url)

	assert.NoError(t, err)
	indirect := r.(*Indirect)
	assert.Equal(t, indirect.URL, url)
}

func TestNotSearchDirectOrSocial(t *testing.T) {
	url := "http://unicorns.ca/"
	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, url, r.(*Indirect).URL)
	assert.Equal(t, r.(*Indirect).Domain, "unicorns.ca")
}

func TestSearchSimple(t *testing.T) {
	r, err := Parse("http://ca.search.yahoo.com/search?p=hello")
	assert.NoError(t, err)
	switch r := r.(type) {
	case *Search:
		assert.Equal(t, r.Label, "Yahoo!")
		assert.Equal(t, r.Domain, "ca.search.yahoo.com")
		assert.Equal(t, r.Query, "hello")
	default:
		assert.Fail(t, fmt.Sprintf("Wrong referrer result: %+v", r))
	}
}

func TestSearchSimpleWithQueryInFragment(t *testing.T) {
	r, err := Parse("http://ca.search.yahoo.com/search#p=hello")
	assert.NoError(t, err)
	switch r := r.(type) {
	case *Search:
		assert.Equal(t, r.Label, "Yahoo!")
		assert.Equal(t, r.Domain, "ca.search.yahoo.com")
		assert.Equal(t, r.Query, "hello")
	default:
		assert.Fail(t, fmt.Sprintf("Wrong referrer result: %+v", r))
	}
}

func TestSearchBingNotLive(t *testing.T) {
	r, err := Parse("http://bing.com/?q=blargh")
	assert.NoError(t, err)
	switch r := r.(type) {
	case *Search:
		assert.Equal(t, r.Label, "Bing")
		assert.Equal(t, r.Domain, "bing.com")
		assert.Equal(t, r.Query, "blargh")
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
	assert.Equal(t, engine.Label, "Yahoo!")
	assert.Equal(t, engine.Domain, "ca.search.yahoo.com")
	assert.True(t, strings.Contains(engine.Query, "\u00F8"))
	assert.Equal(t, engine.Query, "vinduespudsning myshopify rengøring mkobetic")
}

func TestSearchWithExplicitPlus(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+JOKAPOLAR+"11+%2B+11"+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo!")
	assert.Equal(t, engine.Domain, "ca.search.yahoo.com")
	assert.True(t, strings.Contains(engine.Query, "11 + 11"))
	assert.Equal(t, engine.Query, `vinduespudsning JOKAPOLAR "11 + 11" mkobetic`)
}

func TestSearchWithNonAscii(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo!")
	assert.Equal(t, engine.Domain, "ca.search.yahoo.com")
	assert.True(t, strings.Contains(engine.Query, "rengøring"))
	assert.Equal(t, engine.Query, `vinduespudsning myshopify rengøring mkobetic`)
}

func TestSearchWithCyrillics(t *testing.T) {
	url := `http://www.yandex.com/yandsearch?text=%D0%B1%D0%BE%D1%82%D0%B8%D0%BD%D0%BA%D0%B8%20packer-shoes&lr=87&msid=22868.18811.1382712652.60127&noreask=1`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yandex")
	assert.Equal(t, engine.Domain, "www.yandex.com")
	assert.True(t, strings.Contains(engine.Query, "ботинки"))
	assert.Equal(t, engine.Query, `ботинки packer-shoes`)
}

func TestSearchSiteWithEmptyQuery(t *testing.T) {
	url := `https://www.yahoo.com?p=&sa=t&rct=j&p=&esrc=s&source=web&cd=1&ved=0CDkQFjAA&url=http%3A%2F%2Fwww.yellowfashion.in%2F&ei=aZCPUtXmLcGQrQepkIHACA&usg=AFQjCNE-R5-7CENi9oqYe4vG-0g0E7nCSQ&bvm=bv.56988011,d.bmk`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo!")
	assert.Equal(t, engine.Domain, "www.yahoo.com")
	assert.Equal(t, engine.Query, "")
}

func TestSearchSiteGoogleHttpsWithNoParams(t *testing.T) {
	url := `https://www.google.com/`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Google")
	assert.Equal(t, engine.Domain, "www.google.com")
	assert.Equal(t, engine.Query, "")
}

func TestSearchSiteGoogleWithQuery(t *testing.T) {
	url := `https://www.google.co.in/url?sa=t&rct=j&q=test&esrc=s&source=web&cd=1&ved=0CDkQFjAA&url=http%3A%2F%2Fwww.yellowfashion.in%2F&ei=aZCPUtXmLcGQrQepkIHACA&usg=AFQjCNE-R5-7CENi9oqYe4vG-0g0E7nCSQ&bvm=bv.56988011,d.bmk`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Google")
	assert.Equal(t, engine.Domain, "www.google.co.in")
	assert.Equal(t, engine.Query, "test")
}

func TestSearchSiteDuckDuckGoSecured(t *testing.T) {
	url := `http://r.duckduckgo.com/l/?kh=-1&amp;uddg=http%3A%2F%2Fwww.shopify.com%2F`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "DuckDuckGo")
	assert.Equal(t, engine.Domain, "r.duckduckgo.com")
	assert.Equal(t, engine.Query, "")
}

func TestSearchSiteYahooSecured(t *testing.T) {
	url := `http://r.search.yahoo.com/_ylt=A0LEV0H.uiNTzEoA4UjBGOd_;_ylu=X3oDMTByMG04Z2o2BHNlYwNzcgRwb3MDMQRjb2xvA2JmMQR2dGlkAw--/RV=1/RE=1394936959/RO=10/RU=http%3a%2f%2fwww.shopify.com%2f/RS=^ADA0.VXK1194TBSbZf.fSErwtMSJrM-`

	r, err := Parse(url)
	assert.NoError(t, err)

	engine := r.(*Search)
	assert.Equal(t, engine.Label, "Yahoo!")
	assert.Equal(t, engine.Domain, "r.search.yahoo.com")
	assert.Equal(t, engine.Query, "")
}

func TestDirectSimple(t *testing.T) {
	url := "http://example.com"

	r, err := ParseWithDirect(url, "example.com", "sample.com")
	assert.NoError(t, err)

	direct := r.(*Direct)
	assert.NotNil(t, direct)
	assert.Equal(t, direct.URL, url)
	assert.Equal(t, direct.Domain, "example.com")
}

func TestSocialSimple(t *testing.T) {
	url := "https://twitter.com/snormore/status/391149968360103936"

	r, err := Parse(url)
	assert.NoError(t, err)

	social := r.(*Social)
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Twitter")
	assert.Equal(t, social.Domain, "twitter.com")
}

func TestSocialLanguageDomain(t *testing.T) {
	url := "http://es.reddit.com/r/foo"

	r, err := Parse(url)
	assert.NoError(t, err)

	social := r.(*Social)
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Reddit")
	assert.Equal(t, social.Domain, "reddit.com")
}

func TestSocialGooglePlus(t *testing.T) {
	url := "http://plus.url.google.com/url?sa=z&n=1394219098538&url=http%3A%2F%2Fjoe.blogspot.ca&usg=jo2tEVIcI5Wh-6t--v-1ODEeGG8."

	r, err := Parse(url)
	assert.NoError(t, err)

	social := r.(*Social)
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Google+")
	assert.Equal(t, social.Domain, "plus.url.google.com")
}

func TestSocialWithPrefixWWW(t *testing.T) {
	url := "https://www.facebook.com/"

	r, err := Parse(url)
	assert.NoError(t, err)

	social, ok := r.(*Social)
	if !ok {
		assert.Fail(t, "Expected Social", "Instead got %#v", r)
		return
	}
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Facebook")
	assert.Equal(t, social.Domain, "facebook.com")
}

func TestSocialWithPrefixM(t *testing.T) {
	url := "https://m.facebook.com/"

	r, err := Parse(url)
	assert.NoError(t, err)

	social, ok := r.(*Social)
	if !ok {
		assert.Fail(t, "Expected Social", "Instead got %#v", r)
		return
	}
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Facebook")
	assert.Equal(t, social.Domain, "facebook.com")
}

func TestSocialWithPrefixL(t *testing.T) {
	url := "https://l.facebook.com/"

	r, err := Parse(url)
	assert.NoError(t, err)

	social, ok := r.(*Social)
	if !ok {
		assert.Fail(t, "Expected Social", "Instead got %#v", r)
		return
	}
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Facebook")
	assert.Equal(t, social.Domain, "facebook.com")
}

func TestSocialWithPrefixLM(t *testing.T) {
	url := "https://lm.facebook.com/"

	r, err := Parse(url)
	assert.NoError(t, err)

	social, ok := r.(*Social)
	if !ok {
		assert.Fail(t, "Expected Social", "Instead got %#v", r)
		return
	}
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Facebook")
	assert.Equal(t, social.Domain, "facebook.com")
}

func TestEmailSimple(t *testing.T) {
	url := "https://mail.google.com/9aifaufasodf8usafd"

	r, err := Parse(url)
	assert.NoError(t, err)

	email := r.(*Email)
	assert.NotNil(t, email)
	assert.Equal(t, email.Label, "Gmail")
	assert.Equal(t, email.Domain, "mail.google.com")
}

func ExampleParseWithDirect() {
	r, err := ParseWithDirect("http://mysite2.com/products/ties", "mysite1.com", "mysite2.com")
	if err != nil {
		panic(err)
	}
	direct, ok := r.(*Direct)
	if !ok {
		panic("Didn't get a Direct")
	}
	fmt.Printf("Direct %s\n", direct.Domain)
	// Output:
	// Direct mysite2.com
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
			fmt.Printf("Indirect: %s\n", r.URL)
		}
	}
	// Output:
	// Search Yahoo!: hello
	// Social Twitter
	// Indirect: http://yoursite.com/links
}
