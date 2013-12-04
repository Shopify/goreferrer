package referrer

import (
	"encoding/json"
	"io/ioutil"
)

const (
	DataDir       = "./data"
	RulesFilename = "referrers.json"
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

func readRules(rulesPath string) (map[string]SearchRule, map[string]SocialRule, map[string]EmailRule, error) {
	rulesJson, err := ioutil.ReadFile(rulesPath)
	if err != nil {
		return nil, nil, nil, err
	}
	rules := make(map[string]interface{})
	if err = json.Unmarshal(rulesJson, &rules); err != nil {
		return nil, nil, nil, err
	}
	return mappedSearchRules(rules["search"].(map[string]interface{})), mappedSocialRules(rules["social"].(map[string]interface{})), mappedEmailRules(rules["email"].(map[string]interface{})), nil
}

func mappedSearchRules(rawRules map[string]interface{}) map[string]SearchRule {
	mappedRules := make(map[string]SearchRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule.(map[string]interface{})["domains"].([]interface{}) {
			rule := SearchRule{Label: label, Domain: domain.(string)}
			rawParams := rawRule.(map[string]interface{})["parameters"].([]interface{})
			params := make([]string, len(rawParams))
			for _, param := range rawParams {
				params = append(params, param.(string))
			}
			rule.Parameters = params
			mappedRules[rule.Domain] = rule
		}
	}
	return mappedRules
}

func mappedSocialRules(rawRules map[string]interface{}) map[string]SocialRule {
	mappedRules := make(map[string]SocialRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule.(map[string]interface{})["domains"].([]interface{}) {
			mappedRules[domain.(string)] = SocialRule{Label: label, Domain: domain.(string)}
		}
	}
	return mappedRules
}

func mappedEmailRules(rawRules map[string]interface{}) map[string]EmailRule {
	mappedRules := make(map[string]EmailRule)
	for label, rawRule := range rawRules {
		for _, domain := range rawRule.(map[string]interface{})["domains"].([]interface{}) {
			mappedRules[domain.(string)] = EmailRule{Label: label, Domain: domain.(string)}
		}
	}
	return mappedRules
}
