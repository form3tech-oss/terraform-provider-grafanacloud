package grafana

import (
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
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
		SetTimeout(30 * time.Second)

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
