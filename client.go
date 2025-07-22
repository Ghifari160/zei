package zei

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultUserAgent = "Zei/0.1"

// ClientInterface is the interface common to [net/http.Client] and Client.
// Utilizing this interface can assist in migrating existing code from direct usage of
// [net/http.Client] to Client.
type ClientInterface interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (resp *http.Response, err error)
	// Get issues an HTTP GET request to the specified URL.
	Get(url string) (resp *http.Response, err error)
	// Head issues an HTTP HEAD request to the specified URL.
	Head(url string) (resp *http.Response, err error)
	// Post issues an HTTP POST request to the specified URL.
	// If the provided body is an [io.Closer], it is closed after the request.
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	// PostForm issues an HTTP POST request to the specified URL, with data's keys and values
	// URL-encoded as the request body.
	// The Content-Type header is set to application/x-www-form-urlencoded.
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

// A Client is an HTTP client.
// It is built on top of [net/http.Client].
type Client struct {
	conf   *Config
	client *http.Client
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.prepareClient()
	c.setHeaders(req)
	return c.client.Do(req)
}

// Get issues an HTTP GET request to the specified URL.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Head issues an HTTP HEAD request to the specified URL.
func (c *Client) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post issues an HTTP POST request to the specified URL.
// If the provided body is an [io.Closer], it is closed after the request.
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// PostForm issues an HTTP POST request to the specified URL, with data's keys and values
// URL-encoded as the request body.
// The Content-Type header is set to application/x-www-form-urlencoded.
func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// prepareClient configures the underlying [net/http.Client] for the upcoming request.
func (c *Client) prepareClient() {
	c.client.Transport = c.conf.Transport
	c.client.CheckRedirect = c.conf.CheckRedirect
	c.client.Jar = c.conf.Jar
	c.client.Timeout = c.conf.Timeout
}

// setHeaders sets the appropriate headers in req.
func (c *Client) setHeaders(req *http.Request) {
	if c.conf.UserAgent != "" {
		req.Header.Set("User-Agent", c.conf.UserAgent)
	} else {
		req.Header.Set("User-Agent", DefaultUserAgent)
	}
	if c.conf.authMode != authNone {
		req.Header.Set("Authorization", c.conf.authValue)
	}
}

// New creates a new Client.
func New(conf *Config) *Client {
	client := http.Client{
		Timeout: conf.Timeout,
	}

	return &Client{conf, &client}
}
