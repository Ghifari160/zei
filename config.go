package zei

import (
	"net/http"
	"strings"
	"time"
)

// Config configures Client.
// Changes to Config are applied before each [net/http.Request], so it is acceptable to temporarily
// modify a configuration value for one request.
type Config struct {
	// UserAgent specifies the value of the User-Agent header to be sent with each request.
	UserAgent string

	// Transport specifies the mechanism by which individual HTTP requests are made.
	// If nil, [net/http.DefaultTransport] is used.
	Transport http.RoundTripper
	// CheckRedirect specifies the policy for handling redirects.
	// If CheckRedirect is not nil, the client calls it before following an HTTP redirect.
	// The arguments req and via are the upcoming request and the requests already made, oldest
	// first.
	// If CheckRedirect returns an error, the Client's Get method returns both the previous
	// [net/http.Response] (with its body closed) and CheckRedirect's error (wrapped in
	// [url.Error]) instead of issuing the [net/http.Request] req.
	// As a special case, if CheckRedirect returns [net/http.ErrUseLastResponse], then the most
	// recent response is returned with its body unclosed, along with nil error.
	//
	// If CheckRedirect is nil, the Client uses its default policy, which is to stop after 10
	// consecutive requests.
	CheckRedirect func(req *http.Request, via []*http.Request) error
	// Jar specifies the cookie jar.
	//
	// The Jar is used to insert relevant cookies into every outbound [net/http.Request] and is
	// updated with the cookie values of every inbound [net/http.Response].
	// The Jar is consulted for every redirect that the Client follows.
	//
	// If Jar is nil, cookies are only sent if they are explicitly set on the [net/http.Request].
	Jar http.CookieJar

	// Timeout specifies a time limit for requests made by the client.
	// The timeout includes connection time, any redirects, and reading the response body.
	// The timer remains running after Get, Head, Post, or DO return and will interrupt reading of
	// the [net/http.Response.Body].
	Timeout time.Duration

	authMode  authMode
	authValue string
}

// SetBasicAuth configures the Client to set the Authorization header to use HTTP Basic
// Authentication with the provided username and password for every [net/http.Request].
//
// With HTTP Basic Authentication, the provided username and password are not encrypted.
// It should generally only be used in an HTTPS request.
//
// The username may not contain a colon.
// Some protocols may impose additional requirements on pre-escaping the username and password.
func (c *Config) SetBasicAuth(username, password string) {
	c.authMode = authBasic
	c.authValue = username + ":" + password
}

// BasicAuth returns the HTTP BasicAuth Authentication username and password the Client is
// configured to send with every [net/http.Request].
func (c Config) BasicAuth() (username, password string, ok bool) {
	if c.authMode == authBasic {
		auths := strings.SplitN(c.authValue, ":", 2)
		if len(auths) > 0 {
			username = auths[0]
		}
		if len(auths) > 1 {
			password = auths[1]
			ok = true
		}
	}
	return
}

// SetBearerAuth configures the Client to set the Authorization header to use Bearer Authentication
// with the provided token for every [net/http.Request].
func (c *Config) SetBearerAuth(token string) {
	c.authMode = authBearer
	c.authValue = "Bearer " + token
}

// BearerAuth returns the Bearer token the Client is configured to send with every
// [net/http.Request].
func (c Config) BearerAuth() (token string, ok bool) {
	if c.authMode == authBearer {
		auths := strings.SplitN(c.authValue, " ", 2)
		if len(auths) > 1 {
			token = auths[1]
			ok = true
		}
	}
	return
}
