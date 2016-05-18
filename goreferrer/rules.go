package goreferrer

import (
	"encoding/json"
	"io"
	"net/url"
)

type ReferrerType int

const (
	Indirect ReferrerType = iota
	Direct
	Email
	Search
	Social
)

func (r ReferrerType) String() string {
	switch r {
	default:
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

type Parser interface {
	Parse(*url.URL) *Referrer
}

type RuleSet []Parser

// LoadJsonRules can be used to load custom definitions of social sites and search engines.
func LoadJsonRules(reader io.Reader) (RuleSet, error) {
	var decoded jsonRules
	if err := json.NewDecoder(reader).Decode(&decoded); err != nil {
		return nil, err
	}

	var rules RuleSet
	rules = append(rules, extractEmailRules(decoded.Email))
	rules = append(rules, extractSearchRules(decoded.Search))
	rules = append(rules, extractSocialRules(decoded.Social))
	return rules, nil
}

func (r RuleSet) Merge(other RuleSet) {
	for k, v := range other {
		r[k] = v
	}
}

func (r RuleSet) Parse(u *url.URL) *Referrer {
	for _, parser := range r {
		if referrer := parser.Parse(u); referrer != nil {
			return referrer
		}
	}

	return &Referrer{
		Type:   Indirect,
		URL:    u.String(),
		Domain: u.Host,
	}
}

type EmailRule struct {
	Label  string
	Domain string
}

type EmailRules map[string]EmailRule

func (e EmailRules) Parse(u *url.URL) *Referrer {
	if rule, exists := e[u.Host]; exists {
		return &Referrer{
			Type:   Email,
			Label:  rule.Label,
			URL:    u.String(),
			Domain: u.Host,
		}
	}

	return nil
}

type SearchRule struct {
	Label      string
	Domain     string
	Parameters []string
}

type SearchRules map[string]SearchRule

func (s SearchRules) Parse(u *url.URL) *Referrer {
	panic("TODO")
}

type SocialRule struct {
	Label  string
	Domain string
}

type SocialRules map[string]SocialRule

func (s SocialRules) Parse(u *url.URL) *Referrer {
	if rule, exists := s[u.Host]; exists {
		return &Referrer{
			Type:   Social,
			Label:  rule.Label,
			URL:    u.String(),
			Domain: u.Host,
		}
	}

	// Fuzzy search for things like es.reddit.com where the 2-letter locale is the subdomain
	if len(u.Host) > 2 && u.Host[2] == '.' {
		slicedHost := u.Host[3:]

		if rule, exists := s[slicedHost]; exists {
			return &Referrer{
				Type:   Social,
				Label:  rule.Label,
				URL:    u.String(),
				Domain: u.Host,
			}
		}
	}

	return nil
}

func extractEmailRules(ruleMap map[string]jsonRule) EmailRules {
	rules := make(EmailRules)
	for label, jsonRule := range ruleMap {
		for _, domain := range jsonRule.Domains {
			rules[domain] = EmailRule{
				Label:  label,
				Domain: domain,
			}
		}
	}
	return rules
}

func extractSearchRules(ruleMap map[string]jsonRule) SearchRules {
	rules := make(SearchRules)
	for label, jsonRule := range ruleMap {
		for _, domain := range jsonRule.Domains {
			rules[domain] = SearchRule{
				Label:      label,
				Domain:     domain,
				Parameters: jsonRule.Parameters,
			}
		}
	}
	return rules
}

func extractSocialRules(ruleMap map[string]jsonRule) SocialRules {
	rules := make(SocialRules)
	for label, jsonRule := range ruleMap {
		for _, domain := range jsonRule.Domains {
			rules[domain] = SocialRule{
				Label:  label,
				Domain: domain,
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
