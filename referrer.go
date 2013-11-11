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

// engines.csv from https://developers.google.com/analytics/devguides/collection/gajs/gaTrackingTraffic
// Updated on 2013-11-06 (mk)
// Format: label:domain:params
const (
	DataDir         = "./data"
	EnginesFilename = "engines.csv"
	SocialsFilename = "socials.csv"
)

var (
	SearchEngines []Search
	Socials       []Social
	once          sync.Once
)

type Indirect struct {
	Url string
}

type Direct struct {
	Indirect
	Domain string
}

type Search struct {
	Indirect
	Label  string
	Domain string
	Params []string
	Query  string
}

type Social struct {
	Indirect
	Label   string
	Domains []string
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

func readSearchEngines(enginesPath string) ([]Search, error) {
	enginesCsv, err := ioutil.ReadFile(enginesPath)
	if err != nil {
		return nil, err
	}
	var engines []Search
	scanner := bufio.NewScanner(strings.NewReader(string(enginesCsv)))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if line != "" {
			tokens := strings.Split(line, ":")
			params := strings.Split(tokens[2], ",")
			engines = append(engines, Search{Label: tokens[0], Domain: tokens[1], Params: params})
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
			socials = append(socials, Social{Label: tokens[0], Domains: domains})
		}
	}
	return socials, nil
}

func ParseEx(url string, directDomains []string) (interface{}, error) {
	refUrl, err := parseUrl(url)
	if err != nil {
		return nil, err
	}

	if directDomains != nil {
		direct, err := parseDirect(refUrl, directDomains)
		if err != nil {
			return nil, err
		}
		if direct != nil {
			direct.Url = url
			return direct, nil
		}
	}

	social, err := parseSocial(refUrl)
	if err != nil {
		return nil, err
	}
	if social != nil {
		social.Url = url
		return social, nil
	}

	engine, err := parseSearch(refUrl)
	if err != nil {
		return nil, err
	}
	if engine != nil {
		engine.Url = url
		return engine, nil
	}

	return &Indirect{url}, nil
}

func Parse(url string) (interface{}, error) {
	return ParseEx(url, nil)
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
		for _, domain := range social.Domains {
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
	for _, engine := range SearchEngines {
		for _, hostPart := range hostParts {
			if hostPart == engine.Domain {
				for _, param := range engine.Params {
					if search, ok := query[param]; ok {
						return &Search{Label: engine.Label, Query: search[0]}, nil
					}
				}
			}
		}
	}
	return nil, nil
}
