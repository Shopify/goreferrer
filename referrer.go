// Package referrer analyzes and classifies different kinds of referrer URLs (search, social, ...).
package referrer

import (
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	SearchRules   map[string]SearchRule // domain mapping of known rules for search engines.
	SocialRules   map[string]SocialRule // domain mapping of known rules for social sites.
	EmailRules    map[string]EmailRule  // domain mapping of known rules for email sites.
	SearchEngines map[string]SearchRule // list of search engines used for fuzzy matching
	once          sync.Once
)

// Indirect is a referrer that doesn't match any of the other referrer types.
type Indirect struct {
	URL    string // original referrer URL
	Domain string // domain of original referrer URL
}

// Direct is an internal referrer.
// It can only be obtained by calling the extended ParseWithDirect()
type Direct struct {
	URL    string // direct referrer URL
	Domain string // domain of direct referrer URL
}

// Search is a referrer from a set of well known search engines as defined by Google Analytics.
// https://developers.google.com/analytics/devguides/collection/gajs/gaTrackingTraffic.
type Search struct {
	URL    string // search engine referrer URL
	Domain string // matched domain of the search engine, e.g. google.com
	Label  string // search engine label, e.g. Google
	Query  string // decoded search query
}

// Social is a referrer from a set of well know social sites.
type Social struct {
	URL    string // social referrer URL
	Domain string // matched domain of the social site, e.g. twitter.com or t.co
	Label  string // social site label, e.g. Twitter
}

// Email is a referrer from a set of well know email sites.
type Email struct {
	URL    string // email referrer URL
	Domain string // matched domain of the email site, e.g. mail.google.com.com
	Label  string // email site label, e.g. Gmail
}

func init() {
	_, filename, _, _ := runtime.Caller(1)
	once.Do(func() {
		rulesPath := filepath.Join(filepath.Dir(filename), filepath.Join(DataDir, RulesFilename))
		err := InitRules(rulesPath)
		if err != nil {
			panic(err)
		}
		enginesPath := filepath.Join(filepath.Dir(filename), filepath.Join(DataDir, EnginesFilename))
		err = InitSearchEngines(enginesPath)
		if err != nil {
			panic(err)
		}
	})
}

// Parse takes a URL string and turns it into one of the supported referrer types.
// It returns an error if the input is not a valid URL input.
func Parse(url string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refURL, nil)
}

// ParseWithDirect is an extended version of Parse that adds Direct to the set of possible results.
// The additional arguments specify the domains that are to be considered "direct".
func ParseWithDirect(url string, directDomains ...string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refURL, directDomains)
}

func parse(u string, refURL *url.URL, directDomains []string) (interface{}, error) {

	// Parse as direct url
	if directDomains != nil {
		if direct := parseDirect(u, refURL, directDomains); direct != nil {
			return direct, nil
		}
	}

	// Parse as email referrer.
	if email := parseEmail(u, refURL); email != nil {
		return email, nil
	}

	// Parse as social referrer.
	if social := parseSocial(u, refURL); social != nil {
		return social, nil
	}

	// Parse as search referrer.
	if engine := parseSearch(u, refURL); engine != nil {
		return engine, nil
	}

	if engine := fuzzyParseSearch(refURL); engine != nil {
		return engine, nil
	}

	// Parse and return as indirect referrer.
	return &Indirect{URL: u, Domain: refURL.Host}, nil
}

func parseURL(u string) (*url.URL, error) {
	refURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return refURL, nil
}

func parseDirect(rawUrl string, u *url.URL, directDomains []string) *Direct {
	for _, host := range directDomains {
		if host == u.Host {
			return &Direct{URL: rawUrl, Domain: u.Host}
		}
	}
	return nil
}

func parseSocial(rawUrl string, u *url.URL) *Social {
	if rule, ok := SocialRules[u.Host]; ok {
		return &Social{URL: rawUrl, Domain: rule.Domain, Label: rule.Label}
	}
	return nil
}

func parseEmail(rawUrl string, u *url.URL) *Email {
	if rule, ok := EmailRules[u.Host]; ok {
		return &Email{URL: rawUrl, Domain: rule.Domain, Label: rule.Label}
	}
	return nil
}

func parseSearch(rawUrl string, u *url.URL) *Search {
	query := u.Query()
	if rule, ok := SearchRules[u.Host]; ok {
		for _, param := range rule.Parameters {
			if query := query.Get(param); query != "" {
				return &Search{URL: rawUrl, Domain: rule.Domain, Label: rule.Label, Query: query}
			}
		}
	}
	return nil
}

func fuzzyParseSearch(u *url.URL) *Search {
	hostParts := strings.Split(u.Host, ".")
	query := u.Query()
	for _, hostPart := range hostParts {
		if engine, present := SearchEngines[hostPart]; present {
			for _, param := range engine.Parameters {
				if search, ok := query[param]; ok && search[0] != "" {
					return &Search{Label: engine.Label, Query: search[0], Domain: u.Host}
				}
			}
		}
	}
	return nil
}
