package linkChecker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

type checker struct {
	fixedLink    string
	breakLinks   []string
	checkedLinks map[string]bool
	domain       string
	mx           sync.Mutex
	duration     time.Duration
	rl           *rate.Limiter
}

func New(link string) *checker {
	fixMissingLinkProtocol(&link)

	return &checker{
		fixedLink:    link,
		domain:       getLinkDomain(link),
		checkedLinks: map[string]bool{},
		breakLinks:   []string{},
		rl:           rate.NewLimiter(rate.Every(25*time.Millisecond), 50),
	}
}

func (c *checker) Run(maxDepth int) error {
	start := time.Now()

	c.prepereToCheckLinks([]string{c.fixedLink}, 1, &maxDepth)

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

func (c *checker) prepereToCheckLinks(links []string, depth int, maxDepth *int) {
	var wg sync.WaitGroup
	strCh := make(chan string, len(links))

	for _, link := range links {
		c.fixDomainPrefix(&link)

		// Has the url been checked before
		if !c.isCheckedLinks(link) {
			wg.Add(1)
			go c.checkLinks(link, &strCh, &wg, depth, maxDepth)
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

func (c *checker) checkLinks(lnk string, ch *chan string, wg *sync.WaitGroup, depth int, maxDepth *int) {
	defer wg.Done()

	ctx := context.Background()
	err := c.rl.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		log.Println("rl.Wait err: ", err.Error())
		return
	}

	// Send a request / receive a response.
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Get(lnk)
	if err != nil {
		// timeout or connected host has failed to respond
		log.Println("client.Get err: ", err.Error())
		return
	}

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		if response.StatusCode != 404 {
			fmt.Println("response.StatusCode ", response.StatusCode)
		}
		if response.StatusCode != 429 {
			*ch <- lnk
			return
		}
		fmt.Println("response.StatusCode 429")
		wg.Add(1)
		c.checkLinks(lnk, ch, wg, depth, maxDepth)
		return
	}

	if response.StatusCode == 502 {
		// fmt.Println("response.StatusCode: ", response.StatusCode)
		wg.Add(1)
		c.checkLinks(lnk, ch, wg, depth, maxDepth)
		return
	}
	if response.StatusCode != 200 {
		fmt.Println("response.StatusCode: ", response.StatusCode)
	}

	moreLinks := getLinks(html.NewTokenizer(response.Body))

	// Close it manually. To avoid waiting for the end of the function
	if err := response.Body.Close(); err != nil {
		log.Println("Error in the response.Body.Close(). err: ", err.Error())
		return
	}

	if len(moreLinks) > 0 && depth <= *maxDepth {
		c.prepereToCheckLinks(moreLinks, depth+1, maxDepth)
	}
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
