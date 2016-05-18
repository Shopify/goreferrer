package goreferrer

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
