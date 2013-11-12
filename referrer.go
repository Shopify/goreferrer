// Referrer analyzes and classifies different kinds of referrer URLs (search, social, ...).
package referrer

import (
	"bufio"
	"io/ioutil"
	"net/url"
	"path"
	"runtime"
	"strings"
	"sync"
)

const (
	DataDir         = "./data"
	EnginesFilename = "engines.csv"
	SocialsFilename = "socials.csv"
)

var (
	SearchEngines map[string]Search // list of known search engines
	Socials       []Social          // list of known social sites
	once          sync.Once
)

// Indirect is a referrer that doesn't match any of the other referrer types.
type Indirect struct {
	Url string
}

// Direct is an internal referrer.
// It can only be obtained by calling the extended ParseWithDirect()
type Direct struct {
	Indirect
	Domain string
}

// Search is a referrer from a set of well known search engines as defined by Google Analytics.
// https://developers.google.com/analytics/devguides/collection/gajs/gaTrackingTraffic.
type Search struct {
	Indirect
	Label  string
	Query  string
	domain string
	params []string
}

// Social is a referrer from a set of well know social sites.
type Social struct {
	Indirect
	Label   string
	domains []string
}

func init() {
	_, filename, _, _ := runtime.Caller(1)
	once.Do(func() {
		enginesPath := path.Join(path.Dir(filename), path.Join(DataDir, EnginesFilename))
		socialsPath := path.Join(path.Dir(filename), path.Join(DataDir, SocialsFilename))
		err := Init(enginesPath, socialsPath)
		if err != nil {
			panic(err)
		}
	})
}

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
		if line != "" {
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
		if line != "" {
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
	refUrl, err := parseUrl(url)
	if err != nil {
		return nil, err
	}
	return parse(url, refUrl)
}

// ParseWithDirect is an extended version of Parse that adds Direct to the set of possible results.
// The additional arguments specify the domains that are to be considered "direct".
func ParseWithDirect(url string, directDomains ...string) (interface{}, error) {
	refUrl, err := parseUrl(url)
	if err != nil {
		return nil, err
	}
	return parseWithDirect(url, refUrl, directDomains)
}

func parseWithDirect(u string, refUrl *url.URL, directDomains []string) (interface{}, error) {
	if directDomains != nil {
		direct, err := parseDirect(refUrl, directDomains)
		if err != nil {
			return nil, err
		}
		if direct != nil {
			direct.Url = u
			return direct, nil
		}
	}
	return parse(u, refUrl)
}

func parse(u string, refUrl *url.URL) (interface{}, error) {
	social, err := parseSocial(refUrl)
	if err != nil {
		return nil, err
	}
	if social != nil {
		social.Url = u
		return social, nil
	}
	engine, err := parseSearch(refUrl)
	if err != nil {
		return nil, err
	}
	if engine != nil {
		engine.Url = u
		return engine, nil
	}
	return &Indirect{u}, nil
}

func parseUrl(u string) (*url.URL, error) {
	refUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return refUrl, nil
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
				if search, ok := query[param]; ok {
					return &Search{Label: engine.Label, Query: search[0]}, nil
				}
			}
		}
	}
	return nil, nil
}
