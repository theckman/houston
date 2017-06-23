package twilio

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/theckman/houston/twilio/util"
)

// TwilioAPIBase is the base of the Twilio API. Resources will be appended to
// the end of this to make requests to the Twilio API.
const TwilioAPIBase = "https://api.twilio.com/2010-04-01/Accounts"

// Version is the version string of this package. This is used, primary, in
// UserAgent generation.
const Version = "0.0.1"

var userAgent = fmt.Sprintf(
	"Twilio/%s (github.com/theckman/houston/twilio) Go-http-client/1.1 (%s; %s %s)",
	Version, runtime.Version(), runtime.GOOS, runtime.GOARCH,
)

// Client is the struct representing a Twilio client.
type Client struct {
	SID        string
	Secret     string
	HTTPClient HTTPClientInterface
	BaseURL    string
}

// New is a function that takes a sid and secret and returns a *Client. The sid
// and secret value depend on whether you are using your account's Master Keys
// (AccountSid and AuthToken) or generated API keys (API Key SID, API Key
// Secret).
func New(sid, secret string) (*Client, error) {
	if len(sid) == 0 {
		return nil, errors.New("sid cannot be zero length")
	}

	if len(secret) == 0 {
		return nil, errors.New("secret cannot be zero length")
	}

	return &Client{
		SID:        sid,
		Secret:     secret,
		HTTPClient: util.DefaultPooledClient(),
		BaseURL:    TwilioAPIBase,
	}, nil
}

// format takes a resource and ensures it meets the format we expect
// -- which is to say it prepends a '/' if one is not found to ensure
// resource paths resemble "/Path/To/Resource".
func formatResource(resource string) string {
	if len(resource) == 0 {
		return ""
	}

	if strings.Index(resource, "/") != 0 {
		return "/" + resource
	}

	return resource
}

func formatValues(values url.Values) string {
	if len(values) == 0 {
		return ""
	}

	return "?" + values.Encode()
}

func newRequest(client *Client, method, resource string, values url.Values) (*http.Request, error) {
	if client == nil {
		return nil, errors.New("*Client cannot be nil")
	}

	var r *http.Request
	var err error

	switch method {
	case "GET":
		urlStr := fmt.Sprintf(
			"%s/%s%s.json%s",
			client.BaseURL, client.SID,
			formatResource(resource),
			formatValues(values),
		)

		r, err = http.NewRequest(method, urlStr, nil)

		if err != nil {
			return nil, err
		}

		//r.URL.RawQuery = values.Encode()
	case "POST":
		urlStr := fmt.Sprintf(
			"%s/%s%s.json",
			client.BaseURL, client.SID,
			formatResource(resource),
		)

		r, err = http.NewRequest(method, urlStr, strings.NewReader(values.Encode()))

		if err != nil {
			return nil, err
		}
	}

	r.Header.Set("User-Agent", userAgent)
	r.SetBasicAuth(client.SID, client.Secret)

	return r, nil
}

// TestFunc is for development purposes.
//
// TODO(heckman): remove this function before any significant release.
func (c *Client) TestFunc() (*http.Response, error) {
	return c.get("", nil)
}

func (c *Client) get(resource string, params url.Values) (*http.Response, error) {
	req, err := newRequest(c, "GET", resource, params)

	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) post(resource string, formData url.Values) (*http.Response, error) {
	return nil, nil
}
