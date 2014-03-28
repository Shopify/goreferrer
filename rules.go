package referrer

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"strings"
)

const (
	DataDir           = "./data"
	RulesFilename     = "referrers.json"
	EnginesFilename   = "search.csv"
	ParameterWildcard = "*"
)

type SearchRule struct {
	Label      string
	Domain     string
	Parameters []string
}

type SocialRule struct {
	Label  string
	Domain string
}

type EmailRule struct {
	Label  string
	Domain string
}

// RuleSet maps the JSON structure in the file
type RuleSet map[string]map[string][]string

// InitRules can be used to load custom definitions of social sites and search engines
func InitRules(rulesPath string) error {
	rulesJson, err := ioutil.ReadFile(rulesPath)
	if err != nil {
		return err
	}
	rules := make(map[string]RuleSet)
	if err := json.Unmarshal(rulesJson, &rules); err != nil {
		return err
	}
	SearchRules = mappedSearchRules(rules["search"])
	SocialRules = mappedSocialRules(rules["social"])
	EmailRules = mappedEmailRules(rules["email"])
	return nil
}

// InitSearchEngines can be used to load custom definitions of search engines for fuzzy matching
func InitSearchEngines(enginesPath string) error {
	var err error
	SearchEngines, err = readSearchEngines(enginesPath)
	return err
}

func mappedSearchRules(rawRules RuleSet) map[string]SearchRule {
	mappedRules := make(map[string]SearchRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule["domains"] {
			rule := SearchRule{Label: label, Domain: domain}
			rawParams := rawRule["parameters"]
			params := make([]string, len(rawParams))
			for _, param := range rawParams {
				params = append(params, param)
			}
			rule.Parameters = params
			mappedRules[rule.Domain] = rule
		}
	}
	return mappedRules
}

func mappedSocialRules(rawRules RuleSet) map[string]SocialRule {
	mappedRules := make(map[string]SocialRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule["domains"] {
			mappedRules[domain] = SocialRule{Label: label, Domain: domain}
			for _, prefix := range []string{"www.", "m."} {
				variation := prefix + domain
				mappedRules[variation] = SocialRule{Label: label, Domain: domain}
			}
		}
	}
	return mappedRules
}

func mappedEmailRules(rawRules RuleSet) map[string]EmailRule {
	mappedRules := make(map[string]EmailRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule["domains"] {
			mappedRules[domain] = EmailRule{Label: label, Domain: domain}
		}
	}
	return mappedRules
}

func readSearchEngines(enginesPath string) (map[string]SearchRule, error) {
	enginesCsv, err := ioutil.ReadFile(enginesPath)
	if err != nil {
		return nil, err
	}
	engines := make(map[string]SearchRule)
	scanner := bufio.NewScanner(strings.NewReader(string(enginesCsv)))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\r\t")
		if line != "" && line[0] != '#' {
			tokens := strings.Split(line, ":")
			params := strings.Split(tokens[2], ",")
			engines[tokens[1]] = SearchRule{Label: tokens[0], Domain: tokens[1], Parameters: params}
		}
	}
	return engines, nil
}
