// Package referrer analyzes and classifies different kinds of referrer URLs (search, social, ...).
package referrer

import (
	"net/url"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	SearchRules map[string]SearchRule // domain mapping of known rules for search engines.
	SocialRules map[string]SocialRule // domain mapping of known rules for social sites.
	EmailRules  map[string]EmailRule  // domain mapping of known rules for email sites.
	once        sync.Once
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
	Url    string // email referrer URL
	Domain string // matched domain of the email site, e.g. mail.google.com.com
	Label  string // email site label, e.g. Gmail
}

func init() {
	_, filename, _, _ := runtime.Caller(1)
	once.Do(func() {
		rulesPath := path.Join(path.Dir(filename), path.Join(DataDir, RulesFilename))
		err := Init(rulesPath)
		if err != nil {
			panic(err)
		}
	})
}

// Init can be used to load custom definitions of social sites and search engines
func Init(rulesPath string) error {
	var err error
	SearchRules, SocialRules, EmailRules, err = readRules(rulesPath)
	return err
}

// Parse takes a URL string and turns it into one of the supported referrer types.
// It returns an error if the input is not a valid URL input.
func Parse(url string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refUrl, nil)
}

// ParseWithDirect is an extended version of Parse that adds Direct to the set of possible results.
// The additional arguments specify the domains that are to be considered "direct".
func ParseWithDirect(url string, directDomains ...string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refUrl, directDomains)
}

func parse(u string, refUrl *url.URL, directDomains []string) (interface{}, error) {

	// Parse as direct url
	if directDomains != nil {
		direct, err := parseDirect(u, refUrl, directDomains)
		if err != nil {
			return nil, err
		}
		if direct != nil {
			return direct, nil
		}
	}

	// Parse as email referrer.
	email, err := parseEmail(u, refUrl)
	if err != nil {
		return nil, err
	}
	if email != nil {
		return email, nil
	}

	// Parse as social referrer.
	social, err := parseSocial(u, refUrl)
	if err != nil {
		return nil, err
	}
	if social != nil {
		return social, nil
	}

	// Parse as search referrer.
	engine, err := parseSearch(u, refUrl)
	if err != nil {
		return nil, err
	}
	if engine != nil {
		return engine, nil
	}

	// Parse and return as indirect referrer.
	return &Indirect{Url: u, Domain: refUrl.Host}, nil
}

func parseURL(u string) (*url.URL, error) {
	refURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return refURL, nil
}

func parseDirect(rawUrl string, u *url.URL, directDomains []string) (*Direct, error) {
	for _, host := range directDomains {
		if host == u.Host {
			return &Direct{Url: rawUrl, Domain: u.Host}, nil
		}
	}
	return nil, nil
}

func parseSocial(rawUrl string, u *url.URL) (*Social, error) {
	if rule, ok := SocialRules[u.Host]; ok {
		return &Social{Url: rawUrl, Domain: rule.Domain, Label: rule.Label}, nil
	}
	return nil, nil
}

func parseEmail(rawUrl string, u *url.URL) (*Email, error) {
	if rule, ok := EmailRules[u.Host]; ok {
		return &Email{Url: rawUrl, Domain: rule.Domain, Label: rule.Label}, nil
	}
	return nil, nil
}

func parseSearch(rawUrl string, u *url.URL) (*Search, error) {
	query := u.Query()
	if rule, ok := SearchRules[u.Host]; ok {
		for _, param := range rule.Parameters {
			if query := query.Get(param); query != "" {
				return &Search{Url: rawUrl, Domain: rule.Domain, Label: rule.Label, Query: query}, nil
			}
		}
	}
	return nil, nil
}
