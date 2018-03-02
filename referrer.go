package goreferrer

type Referrer struct {
	Type       ReferrerType
	Label      string
	URL        string
	Subdomain  string
	Domain     string
	Tld        string
	Path       string
	Query      string
	GoogleType GoogleSearchType
	Channel    MarketingChannel
}

func (r *Referrer) RegisteredDomain() string {
	if r.Domain != "" && r.Tld != "" {
		return r.Domain + "." + r.Tld
	}

	return ""
}

func (r *Referrer) Host() string {
	if r.Subdomain != "" {
		return r.Subdomain + "." + r.RegisteredDomain()
	}

	return r.RegisteredDomain()
}

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

type GoogleSearchType int

const (
	NotGoogleSearch GoogleSearchType = iota
	OrganicSearch
	Adwords
)

func (g GoogleSearchType) String() string {
	switch g {
	default:
		return "not google search"
	case OrganicSearch:
		return "organic google search"
	case Adwords:
		return "google adwords referrer"
	}
}

type MarketingChannel int

const (
	DirectChannel MarketingChannel = iota
	OrganicSearchChannel
	PaidSearchChannel
	SocialChannel
	EmailChannel
	OrganicReferralChannel
	DisplayChannel
	PaidReferralChannel
)

func (r MarketingChannel) String() string {
	switch r {
	default:
		return "unknown"
	case DirectChannel:
		return "direct"
	case OrganicSearchChannel:
		return "organic_search"
	case PaidSearchChannel:
		return "paid_search"
	case SocialChannel:
		return "social"
	case EmailChannel:
		return "email"
	case OrganicReferralChannel:
		return "organic_referral"
	case DisplayChannel:
		return "display"
	case PaidReferralChannel:
		return "paid_referral"
	}
}
