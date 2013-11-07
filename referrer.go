package referrer

import (
	"bufio"
	"errors"
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
	SearchEngines []SearchEngine
	Socials       []Social
	once          sync.Once
)

type Referrer struct {
	Url string
}

type SearchEngine struct {
	Label  string
	Domain string
	Params []string
	Query  string
}

type Social struct {
	Label   string
	Domains []string
}

type Direct struct {
	Url    string
	Domain string
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

func readSearchEngines(enginesPath string) ([]SearchEngine, error) {
	enginesCsv, err := ioutil.ReadFile(enginesPath)
	if err != nil {
		return nil, err
	}
	var engines []SearchEngine
	scanner := bufio.NewScanner(strings.NewReader(string(enginesCsv)))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if line != "" {
			tokens := strings.Split(line, ":")
			params := strings.Split(tokens[2], ",")
			engines = append(engines, SearchEngine{Label: tokens[0], Domain: tokens[1], Params: params})
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

func NewReferrer(url string) *Referrer {
	r := new(Referrer)
	r.Url = url
	return r
}

func parseUrl(u string) (*url.URL, error) {
	refUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if !refUrl.IsAbs() {
		return nil, errors.New("Referrer URL must be absolute")
	}
	return refUrl, nil
}

func (r *Referrer) ParseDirect(directDomains []string) (*Direct, error) {
	// refUrl, err := parseUrl(r.Url)
	// if err != nil {
	// 	return nil, err
	// }
	// TODO: ...
	return nil, nil
}

func (r *Referrer) ParseSocial() (*Social, error) {
	// refUrl, err := parseUrl(r.Url)
	// if err != nil {
	// 	return nil, err
	// }
	// TODO: ...
	return nil, nil
}

func (r *Referrer) ParseSearchEngine() (*SearchEngine, error) {
	refUrl, err := parseUrl(r.Url)
	if err != nil {
		return nil, err
	}

	hostParts := strings.Split(refUrl.Host, ".")
	query := refUrl.Query()
	for _, engine := range SearchEngines {
		for _, hostPart := range hostParts {
			if hostPart == engine.Domain {
				for _, param := range engine.Params {
					if search, ok := query[param]; ok {
						e := new(SearchEngine)
						e.Query = search[0]
						e.Label = engine.Label
						e.Domain = engine.Domain
						e.Params = engine.Params
						return e, nil
					}
				}
			}
		}
	}

	return new(SearchEngine), nil
}

func (r *Referrer) Parse(directDomains []string) (*Direct, *Social, *SearchEngine, error) {
	direct, err := r.ParseDirect(directDomains)
	if err != nil {
		return nil, nil, nil, err
	}

	social, err := r.ParseSocial()
	if err != nil {
		return nil, nil, nil, err
	}

	engine, err := r.ParseSearchEngine()
	if err != nil {
		return nil, nil, nil, err
	}

	return direct, social, engine, nil
}
