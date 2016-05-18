package goreferrer

import (
	"encoding/json"
	"io"
	"path"
	"strings"
)

type ReferrerType int

const (
	Invalid ReferrerType = iota
	Indirect
	Direct
	Email
	Search
	Social
)

func (r ReferrerType) String() string {
	switch r {
	default:
		return "invalid"
	case Indirect:
		return "indirect"
	case Direct:
		return "direct"
	case Email:
		return "email"
	case Search:
		return "search"
	case Social:
		return "social"
	}
}

type Referrer struct {
	Type      ReferrerType
	Label     string
	URL       string
	Host      string
	Subdomain string
	Domain    string
	Tld       string
	Path      string
	Query     string
}

type Rule struct {
	Type       ReferrerType
	Label      string
	Domain     string
	Parameters []string
}

type RuleSet map[string]Rule

func (r RuleSet) Merge(other RuleSet) {
	for k, v := range other {
		r[k] = v
	}
}

func (r RuleSet) Parse(URL string) Referrer {
	URL = strings.Trim(URL, " \t\r\n")
	if URL == "" {
		return Referrer{
			Type: Direct,
		}
	}

	u, ok := parseUrl(URL)
	if !ok {
		return Referrer{
			Type: Invalid,
			URL:  URL,
		}
	}

	return r.parseUrl(u)
}

func (r RuleSet) ParseWithDirect(URL string, domains ...string) Referrer {
	URL = strings.Trim(URL, " \t\r\n")
	if URL == "" {
		return Referrer{
			Type: Direct,
		}
	}

	u, ok := parseUrl(URL)
	if !ok {
		return Referrer{
			Type: Invalid,
			URL:  URL,
		}
	}

	for _, domain := range domains {
		if u.Host == domain {
			return Referrer{
				Type:      Direct,
				URL:       URL,
				Host:      domain,
				Subdomain: u.Subdomain,
				Domain:    u.Domain,
				Tld:       u.Tld,
				Path:      u.Path,
			}
		}
	}

	return r.parseUrl(u)
}

func (r RuleSet) parseUrl(u *Url) Referrer {
	variations := []string{
		path.Join(u.Host, u.Path),
		path.Join(u.RegisteredDomain(), u.Path),
		u.Host,
		u.RegisteredDomain(),
	}

	for _, host := range variations {
		rule, exists := r[host]
		if !exists {
			continue
		}

		var query string
		for _, param := range rule.Parameters {
			query = u.Query().Get(param)
			if query != "" {
				break
			}
		}

		return Referrer{
			Type:      rule.Type,
			Label:     rule.Label,
			URL:       u.String(),
			Host:      u.Host,
			Subdomain: u.Subdomain,
			Domain:    u.Domain,
			Tld:       u.Tld,
			Path:      u.Path,
			Query:     query,
		}
	}

	return Referrer{
		Type:      Indirect,
		URL:       u.String(),
		Host:      u.Host,
		Subdomain: u.Subdomain,
		Domain:    u.Domain,
		Tld:       u.Tld,
		Path:      u.Path,
	}
}

func extractRules(ruleMap map[string]jsonRule, Type ReferrerType) RuleSet {
	rules := make(RuleSet)
	for label, jsonRule := range ruleMap {
		for _, domain := range jsonRule.Domains {
			rules[domain] = Rule{
				Type:       Type,
				Label:      label,
				Domain:     domain,
				Parameters: jsonRule.Parameters,
			}
		}
	}
	return rules
}

type jsonRule struct {
	Domains    []string
	Parameters []string
}

type jsonRules struct {
	Email  map[string]jsonRule
	Search map[string]jsonRule
	Social map[string]jsonRule
}

// LoadJsonRules can be used to load custom definitions of social sites and
// search engines.
func LoadJsonRules(reader io.Reader) (RuleSet, error) {
	var decoded jsonRules
	if err := json.NewDecoder(reader).Decode(&decoded); err != nil {
		return nil, err
	}

	rules := make(RuleSet)
	rules.Merge(extractRules(decoded.Email, Email))
	rules.Merge(extractRules(decoded.Search, Search))
	rules.Merge(extractRules(decoded.Social, Social))
	return rules, nil
}
