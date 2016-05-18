package goreferrer

import (
	"encoding/json"
	"io"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
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
	Type   ReferrerType
	Label  string
	URL    string
	Domain string
	Query  string
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
			URL:  URL,
		}
	}

	u, err := url.Parse(URL)
	if err != nil || u.Host == "" {
		return Referrer{
			Type: Invalid,
			URL:  URL,
		}
	}

	return r.ParseUrl(u)
}

func (r RuleSet) ParseWithDirect(URL string, domains ...string) Referrer {
	URL = strings.Trim(URL, " \t\r\n")
	if URL == "" {
		return Referrer{
			Type: Direct,
			URL:  URL,
		}
	}

	u, err := url.Parse(URL)
	if err != nil || u.Host == "" {
		return Referrer{
			Type: Invalid,
			URL:  URL,
		}
	}

	for _, domain := range domains {
		if u.Host == domain {
			return Referrer{
				Type:   Direct,
				URL:    URL,
				Domain: domain,
			}
		}
	}

	return r.ParseUrl(u)
}

func (r RuleSet) ParseUrl(u *url.URL) Referrer {
	tld, err := publicsuffix.EffectiveTLDPlusOne(u.Host)
	if err != nil {
		return Referrer{
			Type:   Indirect,
			URL:    u.String(),
			Domain: u.Host,
		}
	}

	variations := []string{
		path.Join(u.Host, u.Path),
		path.Join(tld, u.Path),
		u.Host,
		tld,
	}

	for _, host := range variations {
		rule, exists := r[host]
		if !exists {
			break
		}

		var query string
		for _, param := range rule.Parameters {
			query = u.Query().Get(param)
			if query != "" {
				break
			}
		}

		return Referrer{
			Type:   rule.Type,
			Label:  rule.Label,
			URL:    u.String(),
			Domain: u.Host,
			Query:  query,
		}
	}

	return Referrer{
		Type:   Indirect,
		URL:    u.String(),
		Domain: u.Host,
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
