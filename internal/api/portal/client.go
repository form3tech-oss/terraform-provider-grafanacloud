package portal

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	grafanaStarting = "Your instance is starting"
)

type Client struct {
	client *resty.Client
}

type ClientOpt func(*Client)

func NewClient(baseURL, apiKey string, opts ...ClientOpt) (*Client, error) {
	url := baseURL

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	resty := resty.New().
		SetDebug(len(os.Getenv("HTTP_DEBUG")) != 0).
		SetAuthToken(apiKey).
		SetHostURL(url).
		SetTimeout(10 * time.Second).
		SetRetryWaitTime(10 * time.Second).
		SetRetryCount(6).
		AddRetryCondition(canRetry).
		AddRetryHook(logRetry)

	c := &Client{
		client: resty,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func WithUserAgent(userAgent string) ClientOpt {
	return func(c *Client) {
		c.client.SetHeader("User-Agent", userAgent)
	}
}

// We retry for two reasons:
// 1. Grafana Cloud APIs might apply rate limiting to API requests
// 2. Newly created Grafana Cloud Stacks don't accept requests to create Grafana API keys immediately
func canRetry(r *resty.Response, err error) bool {
	return r.StatusCode() == http.StatusTooManyRequests ||
		strings.Contains(r.String(), grafanaStarting)
}

func logRetry(r *resty.Response, err error) {
	if err != nil {
		log.Printf("[WARN] retrying %s to `%s` because of error: %v", r.Request.Method, r.Request.URL, err)
		return
	}

	log.Printf("[WARN] retrying %s to `%s` because of response: %s", r.Request.Method, r.Request.URL, r)
}
