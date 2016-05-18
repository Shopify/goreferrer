package goreferrer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlankReferrerIsDirect(t *testing.T) {
	blank := DefaultRules.Parse("")
	whitespace := DefaultRules.Parse(" \t\n\r")
	assert.Equal(t, Direct, blank.Type)
	assert.Equal(t, Direct, whitespace.Type)
}

func TestUnknownDomainIsIndirect(t *testing.T) {
	actual := DefaultRules.Parse("https://example.org/path?foo=bar#baz")
	expected := Referrer{
		Type:   Indirect,
		URL:    "https://example.org/path?foo=bar#baz",
		Host:   "example.org",
		Domain: "example",
		Tld:    "org",
		Path:   "/path",
	}
	assert.Equal(t, expected, actual)
}

func TestMalformedUrlIsInvalid(t *testing.T) {
	urls := []string{
		"blap",
		"blap blap",
		"http://",
		"/",
	}

	for _, u := range urls {
		if !assert.Equal(t, Invalid, DefaultRules.Parse(u).Type) {
			t.Log(u)
		}
	}
}

func TestMatchOnAllButQueryString(t *testing.T) {
	rules := RuleSet{
		"www.zambo.com/search": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnDomainTldAndPath(t *testing.T) {
	rules := RuleSet{
		"zambo.com/search": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnSubdomainDomainAndTld(t *testing.T) {
	rules := RuleSet{
		"www.zambo.com": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnDomainAndTld(t *testing.T) {
	rules := RuleSet{
		"zambo.com": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestEmailSimple(t *testing.T) {
	actual := DefaultRules.Parse("https://mail.google.com/9aifaufasodf8usafd")
	expected := Referrer{
		Type:      Email,
		Label:     "Gmail",
		URL:       "https://mail.google.com/9aifaufasodf8usafd",
		Host:      "mail.google.com",
		Subdomain: "mail",
		Domain:    "google",
		Tld:       "com",
		Path:      "/9aifaufasodf8usafd",
	}
	assert.Equal(t, expected, actual)
}

func TestSocialSimple(t *testing.T) {
	actual := DefaultRules.Parse("https://twitter.com/snormore/status/391149968360103936")
	expected := Referrer{
		Type:   Social,
		Label:  "Twitter",
		URL:    "https://twitter.com/snormore/status/391149968360103936",
		Host:   "twitter.com",
		Domain: "twitter",
		Tld:    "com",
		Path:   "/snormore/status/391149968360103936",
	}
	assert.Equal(t, expected, actual)
}

func TestSocialSubdomain(t *testing.T) {
	actual := DefaultRules.Parse("https://puppyanimalbarn.tumblr.com")
	expected := Referrer{
		Type:      Social,
		Label:     "Tumblr",
		URL:       "https://puppyanimalbarn.tumblr.com",
		Host:      "puppyanimalbarn.tumblr.com",
		Subdomain: "puppyanimalbarn",
		Domain:    "tumblr",
		Tld:       "com",
	}
	assert.Equal(t, expected, actual)
}

func TestSocialGooglePlus(t *testing.T) {
	actual := DefaultRules.Parse("http://plus.url.google.com/url?sa=z&n=1394219098538&url=http%3A%2F%2Fjoe.blogspot.ca&usg=jo2tEVIcI5Wh-6t--v-1ODEeGG8.")
	expected := Referrer{
		Type:      Social,
		Label:     "Google+",
		URL:       "http://plus.url.google.com/url?sa=z&n=1394219098538&url=http%3A%2F%2Fjoe.blogspot.ca&usg=jo2tEVIcI5Wh-6t--v-1ODEeGG8.",
		Host:      "plus.url.google.com",
		Subdomain: "plus.url",
		Domain:    "google",
		Tld:       "com",
		Path:      "/url",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchSimple(t *testing.T) {
	actual := DefaultRules.Parse("http://search.yahoo.com/search?p=hello")
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       "http://search.yahoo.com/search?p=hello",
		Host:      "search.yahoo.com",
		Subdomain: "search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search",
		Query:     "hello",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchQueryInFragment(t *testing.T) {
	actual := DefaultRules.Parse("http://search.yahoo.com/search#p=hello")
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       "http://search.yahoo.com/search#p=hello",
		Host:      "search.yahoo.com",
		Subdomain: "search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search",
		Query:     "hello",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchQueryWithYahooCountry(t *testing.T) {
	actual := DefaultRules.Parse("http://ca.search.yahoo.com/search?p=hello")
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       "http://ca.search.yahoo.com/search?p=hello",
		Host:      "ca.search.yahoo.com",
		Subdomain: "ca.search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search",
		Query:     "hello",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchQueryWithYahooCountryAndFragment(t *testing.T) {
	actual := DefaultRules.Parse("http://ca.search.yahoo.com/search#p=hello")
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       "http://ca.search.yahoo.com/search#p=hello",
		Host:      "ca.search.yahoo.com",
		Subdomain: "ca.search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search",
		Query:     "hello",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchBindNotLive(t *testing.T) {
	actual := DefaultRules.Parse("http://bing.com/?q=blargh")
	expected := Referrer{
		Type:   Search,
		Label:  "Bing",
		URL:    "http://bing.com/?q=blargh",
		Host:   "bing.com",
		Domain: "bing",
		Tld:    "com",
		Path:   "/",
		Query:  "blargh",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchNonAscii(t *testing.T) {
	t.Skip("Semicolons are reserved")
	actual := DefaultRules.Parse("http://search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1")
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       "http://search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1",
		Host:      "search.yahoo.com",
		Subdomain: "search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search",
		Query:     "vinduespudsning myshopify rengøring mkobetic",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchWithCyrillics(t *testing.T) {
	actual := DefaultRules.Parse("http://www.yandex.com/yandsearch?text=%D0%B1%D0%BE%D1%82%D0%B8%D0%BD%D0%BA%D0%B8%20packer-shoes&lr=87&msid=22868.18811.1382712652.60127&noreask=1")
	expected := Referrer{
		Type:      Search,
		Label:     "Yandex",
		URL:       "http://www.yandex.com/yandsearch?text=%D0%B1%D0%BE%D1%82%D0%B8%D0%BD%D0%BA%D0%B8%20packer-shoes&lr=87&msid=22868.18811.1382712652.60127&noreask=1",
		Host:      "www.yandex.com",
		Subdomain: "www",
		Domain:    "yandex",
		Tld:       "com",
		Path:      "/yandsearch",
		Query:     "ботинки packer-shoes",
	}
	assert.Equal(t, expected, actual)
}

func TestSearchWithExplicitPlus(t *testing.T) {
	actual := DefaultRules.Parse(`http://search.yahoo.com/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+JOKAPOLAR+"11+%2B+11"+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`)
	expected := Referrer{
		Type:      Search,
		Label:     "Yahoo!",
		URL:       `http://search.yahoo.com/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+JOKAPOLAR+"11+%2B+11"+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`,
		Host:      "search.yahoo.com",
		Subdomain: "search",
		Domain:    "yahoo",
		Tld:       "com",
		Path:      "/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-",
		Query:     `vinduespudsning JOKAPOLAR "11 + 11" mkobetic`,
	}
	assert.Equal(t, expected, actual)
}
