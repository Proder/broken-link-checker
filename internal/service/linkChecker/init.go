package linkChecker

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type checker struct {
	fixedLink    string
	breakLinks   []string
	checkedLinks map[string]bool
	domain       string
	mx           sync.Mutex
	duration     time.Duration
}

func New(link string) *checker {
	fixMissingLinkProtocol(&link)

	return &checker{
		fixedLink:    link,
		domain:       getLinkDomain(link),
		checkedLinks: map[string]bool{},
		breakLinks:   []string{},
	}
}

func (c *checker) Run(maxDepth int) error {
	start := time.Now()

	c.checkLinks([]string{c.fixedLink}, 1, &maxDepth)

	c.duration = time.Since(start)

	return nil
}

func (c *checker) addBreakLink(link *[]string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.breakLinks = append(c.breakLinks, *link...)
}

func (c *checker) isCheckedLinks(link string) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.checkedLinks[link] {
		return true
	}

	c.checkedLinks[link] = true
	return false
}

func (c *checker) checkLinks(links []string, depth int, maxDepth *int) {
	var wg sync.WaitGroup
	strCh := make(chan string, len(links))

	for _, link := range links {
		c.fixDomainPrefix(&link)

		// Has the url been checked before
		if !c.isCheckedLinks(link) {

			wg.Add(1)
			go func(lnk string, ch *chan string) {
				defer wg.Done()

				// Send a request / receive a response.
				client := http.Client{
					Timeout: 60 * time.Second,
				}
				response, err := client.Get(lnk)
				if err != nil {
					log.Println("client.Get err: ", err.Error())
					return
				}

				if response.StatusCode >= 400 && response.StatusCode < 500 {
					*ch <- lnk
					return
				}

				moreLinks := getLinks(html.NewTokenizer(response.Body))

				// Close it manually. To avoid waiting for the end of the function
				if err := response.Body.Close(); err != nil {
					log.Println("Error in the response.Body.Close(). err: ", err.Error())
					return
				}

				if len(moreLinks) > 0 && depth <= *maxDepth {
					c.checkLinks(moreLinks, depth+1, maxDepth)
				}
			}(link, &strCh)
		}
	}

	go func() {
		wg.Wait()
		close(strCh)
	}()

	var brLinks []string
	for v := range strCh {
		brLinks = append(brLinks, v)
	}
	c.addBreakLink(&brLinks)
}

func (c *checker) fixDomainPrefix(link *string) {
	if link == nil {
		return
	}

	if !strings.Contains(*link, "://") {
		*link = c.domain + *link
	}
}

func (c *checker) GetBreakLinks() []string {
	return c.breakLinks
}

func (c *checker) GetDuration() string {
	return c.duration.String()
}
