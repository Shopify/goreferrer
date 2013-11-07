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
	url string
}

type SearchEngine struct {
	Label  string
	Domain string
	Params []string
}

type Social struct {
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
	r.url = url
	return r
}

func (r *Referrer) ParseDirect(directDomains []string) error {
	return nil
}

func (r *Referrer) ParseSocial() error {
	return nil
}

func (r *Referrer) ParseSearchEngine() error {
	return nil
}

func (r *Referrer) Parse(directDomains []string) error {
	return nil
}

func analyzeReferrer(urlString string) (map[string]interface{}, error) {
	data := map[string]interface{}{"referrer": urlString}
	referrer, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if !referrer.IsAbs() {
		return nil, errors.New("Referrer URL must be absolute")
	}
	host := strings.Split(referrer.Host, ".")
	query := referrer.Query()
	for _, engine := range SearchEngines {
		for _, name := range host {
			if name == engine.Domain {
				for _, param := range engine.Params {
					if search, ok := query[param]; ok {
						data["kind"] = "s"
						data["engine"] = engine.Label
						data["query"] = search[0]
						return data, nil
					}
				}
			}
		}
	}
	data["kind"] = "r"
	return data, nil
}
