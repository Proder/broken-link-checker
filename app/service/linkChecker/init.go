package linkChecker

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Checker struct {
	breakLinks   []string
	checkedLinks map[string]bool
	domain       string
	mx           sync.Mutex
}

func (c *Checker) Run(link string, maxDepth int) error {
	correctLink := fixProtocolPrefix(link)

	c.domain = getDomain(correctLink)
	c.checkedLinks = make(map[string]bool)

	c.checkLinks([]string{correctLink}, 1, &maxDepth)

	fmt.Printf("Broken links found: %d\n", len(c.breakLinks))

	return nil
}

func (c *Checker) addBreakLink(link *[]string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.breakLinks = append(c.breakLinks, *link...)
}

func (c *Checker) isCheckedLinks(link *string) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	if !c.checkedLinks[*link] {
		c.checkedLinks[*link] = true
		return false
	} else {
		return true
	}
}

func (c *Checker) checkLinks(links []string, depth int, maxDepth *int) {
	var wg sync.WaitGroup
	strCh := make(chan string, len(links))

	for _, link := range links {
		c.fixDomainPrefix(&link)

		// Has the url been checked before
		if !c.isCheckedLinks(&link) {

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

				moreLinks := getLinks(html.NewTokenizer(response.Body))

				// Close it manually. To avoid waiting for the end of the function
				if err := response.Body.Close(); err != nil {
					_ = fmt.Errorf("Error in the response.Body.Close(). err: %s ", err)
					return
				}

				if len(moreLinks) > 0 && depth < *maxDepth {
					c.checkLinks(moreLinks, depth+1, maxDepth)
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
	c.addBreakLink(&brLinks)
}

func (c *Checker) fixDomainPrefix(link *string) *string {
	if !strings.Contains(*link, "://") {
		*link = c.domain + *link
	}

	return link
}
