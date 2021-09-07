package linkChecker

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

type Info struct {
	ErrConn   int    `json:"errConn"`
	ErrServer int    `json:"errServer"`
	Duration  string `json:"duration"`
}

type checker struct {
	fixedLink    string
	breakLinks   []string
	info         Info
	checkedLinks map[string]bool
	domain       string
	maxDepth     int
	mx           sync.Mutex
	rl           *rate.Limiter
}

func New(link string, maxDepth int) *checker {
	fixMissingLinkProtocol(&link)

	return &checker{
		fixedLink:    link,
		domain:       getLinkDomain(link),
		maxDepth:     maxDepth,
		checkedLinks: map[string]bool{},
		breakLinks:   []string{},
		rl:           rate.NewLimiter(rate.Every(25*time.Millisecond), 50),
		info:         Info{},
	}
}

func (c *checker) Run() error {
	start := time.Now()

	c.checkLinks([]string{c.fixedLink}, 1)

	c.info.Duration = time.Since(start).String()

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

func (c *checker) checkLinks(links []string, depth int) {
	var wg sync.WaitGroup
	strCh := make(chan string, len(links))

	for _, link := range links {
		c.fixDomainPrefix(&link)

		// Has the url been checked before
		if !c.isCheckedLinks(link) {
			wg.Add(1)

			ctx := context.Background()
			err := c.rl.Wait(ctx) // This is a blocking call. Honors the rate limit
			if err != nil {
				log.Println("rl.Wait err: ", err.Error())
				return
			}

			go func(lnk string, depth int, ch *chan string) {
				defer wg.Done()

				// Send a request / receive a response.
				client := http.Client{
					Timeout: 60 * time.Second,
				}
				response, err := client.Get(lnk)
				if err != nil {
					// timeout or connected host has failed to respond
					log.Println("client.Get err: ", err.Error())
					c.info.ErrConn++
					return
				}

				if response.StatusCode >= 400 && response.StatusCode < 500 {
					*ch <- lnk
					return
				}

				if response.StatusCode > 500 {
					log.Println("server error. response.StatusCode: ", response.StatusCode)
					c.info.ErrServer++
					return
				}

				moreLinks := getLinks(html.NewTokenizer(response.Body))

				// Close it manually. To avoid waiting for the end of the function
				if err := response.Body.Close(); err != nil {
					log.Println("Error in the response.Body.Close(). err: ", err.Error())
					return
				}

				if len(moreLinks) > 0 && depth <= c.maxDepth {
					c.checkLinks(moreLinks, depth+1)
				}
			}(link, depth, &strCh)
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
	return c.info.Duration
}

func (c *checker) GetInfo() Info {
	return c.info
}
