// Package referrer analyzes and classifies different kinds of referrer URLs (search, social, ...).
package referrer

import (
	"bufio"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	// DataDir contains the CSV files listing recognized search engine and social referrers.
	DataDir = "data"
	// EnginesFilename is the CSV file recognizing search engine referrers.
	EnginesFilename = "engines.csv"
	// SocialsFilename is the CSV file recognizing social referrers.
	SocialsFilename = "socials.csv"
)

var (
	// SearchEngines is a map of known search engines
	SearchEngines map[string]Search
	// Socials is a list of known social sites
	Socials []Social
	once    sync.Once
)

// Indirect is a referrer that doesn't match any of the other referrer types.
type Indirect struct {
	URL string // original referrer URL
}

// Direct is an internal referrer.
// It can only be obtained by calling the extended ParseWithDirect()
type Direct struct {
	Indirect
	Domain string // direct domain that matched the URL
}

// Search is a referrer from a set of well known search engines as defined by Google Analytics.
// https://developers.google.com/analytics/devguides/collection/gajs/gaTrackingTraffic.
type Search struct {
	Indirect
	Label  string // search engine label, e.g Google
	Query  string // decoded search query
	domain string
	params []string
}

// Social is a referrer from a set of well know social sites.
type Social struct {
	Indirect
	Label   string // social site label, e.g. Twitter
	domains []string
}

func init() {
	_, filename, _, _ := runtime.Caller(1)
	once.Do(func() {
		enginesPath := filepath.Join(filepath.Dir(filename), filepath.Join(DataDir, EnginesFilename))
		socialsPath := filepath.Join(filepath.Dir(filename), filepath.Join(DataDir, SocialsFilename))
		err := Init(enginesPath, socialsPath)
		if err != nil {
			panic(err)
		}
	})
}

// Init can be used to load custom definitions of social sites and search engines
func Init(enginesPath string, socialsPath string) error {
	var err error
	SearchEngines, err = readSearchEngines(enginesPath)
	Socials, err = readSocials(socialsPath)
	return err
}

func readSearchEngines(enginesPath string) (map[string]Search, error) {
	enginesCsv, err := ioutil.ReadFile(enginesPath)
	if err != nil {
		return nil, err
	}
	engines := make(map[string]Search)
	scanner := bufio.NewScanner(strings.NewReader(string(enginesCsv)))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if line != "" && line[0] != '#' {
			tokens := strings.Split(line, ":")
			params := strings.Split(tokens[2], ",")
			engines[tokens[1]] = Search{Label: tokens[0], domain: tokens[1], params: params}
		}
	}
	return engines, nil
}

func readSocials(socialsPath string) ([]Social, error) {
	socialsCsv, err := ioutil.ReadFile(socialsPath)
	if err != nil {
		return nil, err
	}
	var socials []Social
	scanner := bufio.NewScanner(strings.NewReader(string(socialsCsv)))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if line != "" && line[0] != '#' {
			tokens := strings.Split(line, ":")
			domains := strings.Split(tokens[1], ",")
			socials = append(socials, Social{Label: tokens[0], domains: domains})
		}
	}
	return socials, nil
}

// Parse takes a URL string and turns it into one of the supported referrer types.
// It returns an error if the input is not a valid URL input.
func Parse(url string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refURL)
}

// ParseWithDirect is an extended version of Parse that adds Direct to the set of possible results.
// The additional arguments specify the domains that are to be considered "direct".
func ParseWithDirect(url string, directDomains ...string) (interface{}, error) {
	refURL, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	return parseWithDirect(url, refURL, directDomains)
}

func parseWithDirect(u string, refURL *url.URL, directDomains []string) (interface{}, error) {
	if directDomains != nil {
		direct, err := parseDirect(refURL, directDomains)
		if err != nil {
			return nil, err
		}
		if direct != nil {
			direct.URL = u
			return direct, nil
		}
	}
	return parse(u, refURL)
}

func parse(u string, refURL *url.URL) (interface{}, error) {
	social, err := parseSocial(refURL)
	if err != nil {
		return nil, err
	}
	if social != nil {
		social.URL = u
		return social, nil
	}
	engine, err := parseSearch(refURL)
	if err != nil {
		return nil, err
	}
	if engine != nil {
		engine.URL = u
		return engine, nil
	}
	return &Indirect{u}, nil
}

func parseURL(u string) (*url.URL, error) {
	refURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return refURL, nil
}

func parseDirect(u *url.URL, directDomains []string) (*Direct, error) {
	for _, host := range directDomains {
		if host == u.Host {
			return &Direct{Domain: host}, nil
		}
	}
	return nil, nil
}

func parseSocial(u *url.URL) (*Social, error) {
	for _, social := range Socials {
		for _, domain := range social.domains {
			if domain == u.Host {
				return &Social{Label: social.Label}, nil
			}
		}
	}
	return nil, nil
}

func parseSearch(u *url.URL) (*Search, error) {
	hostParts := strings.Split(u.Host, ".")
	query := u.Query()
	for _, hostPart := range hostParts {
		if engine, present := SearchEngines[hostPart]; present {
			for _, param := range engine.params {
				if search, ok := query[param]; ok && search[0] != "" {
					return &Search{Label: engine.Label, Query: search[0]}, nil
				}
			}
		}
	}
	return nil, nil
}
