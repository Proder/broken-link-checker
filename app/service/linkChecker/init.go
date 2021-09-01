package linkChecker

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var (
	breakLinks   []string
	checkedLinks = make(map[string]bool)
	domain       string
)

func Run(link string, maxDepth int) error {
	l := fixProtocolPrefix(link)

	s := strings.Split(l, "/")
	domain = s[0] + "//" + s[2]

	checkLinks([]string{l}, 1, &maxDepth)

	fmt.Printf("Broken links found: %d\n", len(breakLinks))
	fmt.Println(breakLinks)
	return nil
}

func fixProtocolPrefix(link string) string {
	if !strings.Contains(link, "://") {
		link = "http://" + link
	}
	return link
}

func fixDomainPrefix(link *string) *string {
	if !strings.Contains(*link, "://") {
		*link = domain + *link
	}

	return link
}

func checkLinks(links []string, depth int, maxDepth *int) {
	for _, link := range links {
		fixDomainPrefix(&link)
		// Has the url been checked before
		if !checkedLinks[link] {
			checkedLinks[link] = true

			// Send a request / receive a response.
			response, err := http.Get(link)
			if err != nil {
				fmt.Println("http.Get err: " + err.Error())

				continue
			}

			if response.StatusCode >= 400 && response.StatusCode < 500 {
				breakLinks = append(breakLinks, link)
			}

			moreLinks := getMoreLinks(html.NewTokenizer(response.Body))
			// Close it manually. To avoid waiting for the end of the function
			if err := response.Body.Close(); err != nil {
				_ = fmt.Errorf("Error in the response.Body.Close(). err: %s ", err)
				continue
			}

			if len(moreLinks) > 0 && depth < *maxDepth {
				depth++
				checkLinks(moreLinks, depth, maxDepth)
			}
		}
	}
}

func getMoreLinks(htmlTokens *html.Tokenizer) (newLink []string) {
	for {
		tt := htmlTokens.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := htmlTokens.Token()
			if t.Data == "a" {
				newLink = append(newLink, t.Attr[0].Val)
			}
		}
	}
}
