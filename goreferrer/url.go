package goreferrer

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Url struct {
	*url.URL
	Subdomain string
	Domain    string
	Tld       string
}

func parseUrl(s string) (*Url, bool) {
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return nil, false
	}

	tld, _ := publicsuffix.PublicSuffix(u.Host)
	if tld == "" || len(u.Host)-len(tld) < 2 {
		return nil, false
	}

	hostWithoutTld := u.Host[:len(u.Host)-len(tld)-1]
	lastDot := strings.LastIndex(hostWithoutTld, ".")
	if lastDot == -1 {
		return &Url{URL: u, Domain: hostWithoutTld, Tld: tld}, true
	}

	return &Url{
		URL:       u,
		Subdomain: hostWithoutTld[:lastDot],
		Domain:    hostWithoutTld[lastDot+1:],
		Tld:       tld,
	}, true
}

func (u *Url) RegisteredDomain() string {
	return u.Domain + "." + u.Tld
}
