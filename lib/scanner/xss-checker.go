package scanner

import (
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

type XssChecker struct{
	Payload []string
	Target string
	Parameters interface{}
}

func (xss *XssChecker) Scan() []Vulnerabilities {
	var xssVuln Vulnerabilities
	var xssVulns []Vulnerabilities
	res, err := http.Get(xss.Target)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	// Parse the HTML of the webpage
	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Use a regular expression to search for potential XSS vulnerabilities
	pattern := `(?i)(<[^>]+\b(src|href|style|on[a-z]+)\s*=\s*['"]?[^'">]*\b(alert|prompt|confirm)\([^>]*>)`
	re := regexp.MustCompile(pattern)
	// Traverse the DOM tree and search for XSS vulnerabilities
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if re.MatchString(n.Data) {
				// Found an XSS vulnerability
				if n.Attr != nil {
					// Check if the vulnerability is stored or reflected
					for _, attr := range n.Attr {
						if attr.Key == "src" || attr.Key == "href" {
							fmt.Println("Possible reflected XSS vulnerability found in attribute:", attr.Val)
						} else if attr.Key == "style" || (attr.Key[:2] == "on" && attr.Key[2:] != "error") {
							fmt.Println("Possible stored XSS vulnerability found in attribute:", attr.Val)
						}
					}
				} else {
					// If there are no attributes, it is likely a DOM-based XSS vulnerability
					fmt.Println("Possible DOM-based XSS vulnerability found")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
