package goreferrer

import (
	"encoding/json"
	"io"
	"net/url"
	"path"
	"strings"
)

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
				Path:      cleanPath(u.Path),
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

		query := getQuery(u.Query(), rule.Parameters)
		if query == "" {
			values, err := url.ParseQuery(u.Fragment)
			if err == nil {
				query = getQuery(values, rule.Parameters)
			}
		}

		ref := Referrer{
			Type:      rule.Type,
			Label:     rule.Label,
			URL:       u.String(),
			Host:      u.Host,
			Subdomain: u.Subdomain,
			Domain:    u.Domain,
			Tld:       u.Tld,
			Path:      cleanPath(u.Path),
			Query:     query,
		}
		ref.GoogleType = googleSearchType(ref)
		return ref
	}

	return Referrer{
		Type:      Indirect,
		Label:     strings.Title(u.Domain),
		URL:       u.String(),
		Host:      u.Host,
		Subdomain: u.Subdomain,
		Domain:    u.Domain,
		Tld:       u.Tld,
		Path:      cleanPath(u.Path),
	}
}

func getQuery(values url.Values, params []string) string {
	for _, param := range params {
		query := values.Get(param)
		if query != "" {
			return query
		}
	}

	return ""
}

func googleSearchType(ref Referrer) GoogleSearchType {
	if ref.Type != Search || !strings.Contains(ref.Label, "Google") {
		return NotGoogleSearch
	}

	if strings.HasPrefix(ref.Path, "/aclk") || strings.HasPrefix(ref.Path, "/pagead/aclk") {
		return Adwords
	}

	return OrganicSearch
}

func cleanPath(path string) string {
	if i := strings.Index(path, ";"); i != -1 {
		return path[:i]
	}
	return path
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
