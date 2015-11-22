package uber

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// Uber API endpoint
	APIUrl string = "https://api.uber.com/v1/%s%s"
)

// Getter defines the behavior for all HTTP Get requests
type Getter interface {
	get(c *Client) error
}

// OAuth parameters
type RequestOptions struct {
	ServerToken    string
	ClientId       string
	ClientSecret   string
	AppName        string
	AuthorizeUrl   string
	AccessTokenUrl string
	BaseUrl        string
}

// Client contains the required OAuth tokens and urls and manages
// the connection to the API. All requests are made via this type
type Client struct {
	Options *RequestOptions
}

// Create returns a new API client
func Create(options *RequestOptions) *Client {
	return &Client{options}
}

// Get formulates an HTTP GET request based on the Uber endpoint type
func (c *Client) Get(getter Getter) error {
	if e := getter.get(c); e != nil {
		return e
	}

	return nil
}

// Send HTTP request to Uber API
func (c *Client) getRequest(endpoint string, params map[string]string) []byte {
	urlParams := "?"
	params["server_token"] = c.Options.ServerToken
	for k, v := range params {
		if len(urlParams) > 1 {
			urlParams += "&"
		}
		urlParams += fmt.Sprintf("%s=%s", k, v)
	}

	url := fmt.Sprintf(APIUrl, endpoint, urlParams)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}
