goreferrer
==========

A Go package that analyzes and classifies different kinds of referrer URLs (search, social, ...).

## Example

```
	import ("github.com:Shopify/goreferrer")
	
	urls := []string{
		"http://ca.search.yahoo.com/search?p=hello",
		"https://twitter.com/jdoe/status/391149968360103936",
		"http://mysite.com/links"
		"http://yoursite.com/links"
	}
	for url := range urls {
		r, err := referrer.ParseWithDirect(url,"mysite.com")
		switch r := r.(type) {
		case *Search: fmt.Printf("Search %s: %s",r.Label, r.Query)
		case *Social: fmt.Printf("Social %s", r.Label)
		case *Direct: fmt.Printf("Direct %s", r.Domain)
		case *Indirect: fmt.Printf("Indirect: %s", r.Url)
		}
	}
```
Result:
```
	Search Yahoo: hello
	Social Twitter
	Direct mysite.com
	Indirect: http://mysite.com/links
```
