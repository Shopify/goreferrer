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
		Domain: "example.org",
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
