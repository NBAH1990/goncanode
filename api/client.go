package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
)

type IClient interface {
	Request(ctx context.Context, method string, url string, data *bytes.Buffer) (result []byte, err error)
}

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
}

func (c *Client) Request(ctx context.Context, method string, url string, data *bytes.Buffer) (result []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, c.resolveTrimmedUrl(url), data)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (c *Client) resolveTrimmedUrl(url string) string {
	p1 := strings.TrimRight(c.BaseUrl, "/")
	p2 := strings.TrimLeft(url, "/")

	if p2 == "" {
		return p1
	}

	return p1 + "/" + p2
}
