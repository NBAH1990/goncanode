package api

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Url string
}

func (c *Client) Request(ctx context.Context, method string, data *bytes.Buffer) (result []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, c.Url, data)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return result, nil
}
