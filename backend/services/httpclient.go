package services

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
	PostJSON(url string, jsonData []byte) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type HttpClientWrapper struct {
	client *http.Client
}

func ConstructHttpClient() *HttpClientWrapper {
	defaultHeaders := http.Header{
		"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"},
		"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Accept-Encoding": []string{"gzip, deflate, br"},
		"Accept-Language": []string{"en-US,en;q=0.9,lt;q=0.8,et;q=0.7,de;q=0.6"},
	}
	return &HttpClientWrapper{
		client: &http.Client{
			Transport: &HeaderRoundTripper{
				Transport: http.DefaultTransport,
				Headers:   defaultHeaders,
			},
			Timeout: 10 * time.Second,
		},
	}
}
func (c *HttpClientWrapper) Get(url string) (*http.Response, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClientWrapper) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	resp, err := c.client.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClientWrapper) PostJSON(url string, jsonData []byte) (*http.Response, error) {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClientWrapper) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
