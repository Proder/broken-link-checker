package linkChecker

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	breakLinks   []string
	checkedLinks = make(map[string]bool)
	domain       string
	mx           sync.Mutex
	rwMx         sync.Mutex
)

func Run(link string, maxDepth int) error {
	l := fixProtocolPrefix(link)

	s := strings.Split(l, "/")
	domain = s[0] + "//" + s[2]

	checkLinks([]string{l}, 1, &maxDepth)

	fmt.Printf("Broken links found: %d\n", len(breakLinks))

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

func addBreakLink(link *[]string) {
	mx.Lock()
	defer mx.Unlock()

	breakLinks = append(breakLinks, *link...)
}

func isCheckedLinks(link string) bool {
	rwMx.Lock()
	defer rwMx.Unlock()

	if !checkedLinks[link] {
		checkedLinks[link] = true
		return false
	} else {
		return true
	}
}

func checkLinks(links []string, depth int, maxDepth *int) {
	var wg sync.WaitGroup
	strCh := make(chan string, len(links))

	for _, link := range links {
		fixDomainPrefix(&link)

		// Has the url been checked before
		if !isCheckedLinks(link) {

			wg.Add(1)
			go func(lnk string, ch *chan string, wg *sync.WaitGroup) {
				defer wg.Done()

				// Send a request / receive a response.
				response, err := http.Get(lnk)
				if err != nil {
					fmt.Println("http.Get err: " + err.Error())

					return
				}

				if response.StatusCode >= 400 && response.StatusCode < 500 {
					*ch <- lnk
					return
				}

				moreLinks := getMoreLinks(html.NewTokenizer(response.Body))
				// Close it manually. To avoid waiting for the end of the function
				if err := response.Body.Close(); err != nil {
					_ = fmt.Errorf("Error in the response.Body.Close(). err: %s ", err)
					return
				}

				if len(moreLinks) > 0 && depth < *maxDepth {
					checkLinks(moreLinks, depth+1, maxDepth)
				}
			}(link, &strCh, &wg)
		}
	}

	wg.Wait()
	close(strCh)

	var brLinks []string
	for v := range strCh {
		brLinks = append(brLinks, v)
	}
	addBreakLink(&brLinks)
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
