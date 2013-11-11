package referrer

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	ExampleDirectDomains = []string{"example.com", "sample.com"}
)

func TestRelativeUrl(t *testing.T) {
	url := `buh`
	r, err := Parse(url)

	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestNotSearchDirectOrSocial(t *testing.T) {
	url := "http://unicorns.ca/"
	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindIndirect)
	assert.Equal(t, url, r.(*Indirect).Url)
}

func TestSearchNonAscii(t *testing.T) {
	url := "http://ca.search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1"
	assert.True(t, strings.Contains(url, "\u00F8"))

	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindSearchEngine)

	engine := r.(*SearchEngine)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "\u00F8"))
	assert.Equal(t, engine.Query, "vinduespudsning myshopify rengøring mkobetic")
}

func TestSearchWithExplicitPlus(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8nVvm5StDIAIxHrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDSjNTOW9rZ2V1eVVMYVp6c1VmRmRMUkdDMkxfbjJsSnV2dFVBQmZyWgRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANDc01MSGlnTVFOS2k2cDRqcUxERzRBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgSk9LQVBPTEFSICIxMSArIDExIiBta29iZXRpYwR0X3N0bXADMTM4Mjk4OTYwMjg3OQR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+JOKAPOLAR+"11+%2B+11"+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindSearchEngine)

	engine := r.(*SearchEngine)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "11 + 11"))
	assert.Equal(t, engine.Query, `vinduespudsning JOKAPOLAR "11 + 11" mkobetic`)
}

func TestSearchWithNonAscii(t *testing.T) {
	url := `http://ca.search.yahoo.com/search;_ylt=A0geu8fBeW5SqVEAZ2vrFAx.;_ylc=X1MDMjExNDcyMTAwMwRfcgMyBGJjawMwbXFjc3RoOHYybjlkJTI2YiUzRDMlMjZzJTNEYWkEY3NyY3B2aWQDWmxUdFhVZ2V1eVVMYVp6c1VmRmRMUXUyMkxfbjJsSnVlY0VBQlhDWQRmcgN5ZnAtdC03MTUEZnIyA3NiLXRvcARncHJpZANVRFRzSGFBUVF0ZUZHZ2hzZ0N3VDNBBG10ZXN0aWQDbnVsbARuX3JzbHQDMARuX3N1Z2cDMARvcmlnaW4DY2Euc2VhcmNoLnlhaG9vLmNvbQRwb3MDMARwcXN0cgMEcHFzdHJsAwRxc3RybAM0NARxdWVyeQN2aW5kdWVzcHVkc25pbmcgbXlzaG9waWZ5IHJlbmf4cmluZyBta29iZXRpYwR0X3N0bXADMTM4Mjk3MjM1NDIzMwR2dGVzdGlkA01TWUNBQzE-?p=vinduespudsning+myshopify+rengøring+mkobetic&fr2=sb-top&fr=yfp-t-715&rd=r1`

	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindSearchEngine)

	engine := r.(*SearchEngine)
	assert.Equal(t, engine.Label, "Yahoo")
	assert.True(t, strings.Contains(engine.Query, "rengøring"))
	assert.Equal(t, engine.Query, `vinduespudsning myshopify rengøring mkobetic`)
}

func TestSearchWithCyrillics(t *testing.T) {
	url := `http://www.yandex.com/yandsearch?text=%D0%B1%D0%BE%D1%82%D0%B8%D0%BD%D0%BA%D0%B8%20packer-shoes&lr=87&msid=22868.18811.1382712652.60127&noreask=1`

	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindSearchEngine)

	engine := r.(*SearchEngine)
	assert.Equal(t, engine.Label, "Yandex")
	assert.True(t, strings.Contains(engine.Query, "ботинки"))
	assert.Equal(t, engine.Query, `ботинки packer-shoes`)
}

func TestDirectSimple(t *testing.T) {
	url := "http://example.com"

	r, err := ParseEx(url, ExampleDirectDomains)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindDirect)

	direct := r.(*Direct)
	assert.NotNil(t, direct)
	assert.Equal(t, direct.Url, url)
	assert.Equal(t, direct.Domain, "example.com")
}

func TestSocialSimple(t *testing.T) {
	url := "https://twitter.com/snormore/status/391149968360103936"

	r, err := Parse(url)
	assert.NoError(t, err)
	assert.Equal(t, r.Kind(), KindSocial)

	social := r.(*Social)
	assert.NotNil(t, social)
	assert.Equal(t, social.Label, "Twitter")
}
