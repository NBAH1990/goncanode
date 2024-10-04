package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestClient_Request_Success(t *testing.T) {
	c := &Client{
		BaseUrl: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if req.Method != "GET" {
						t.Errorf("Expected method GET, got %s", req.Method)
					}
					if req.URL.String() != "https://example.com/path" {
						t.Errorf("Expected URL https://example.com/path, got %s", req.URL.String())
					}
					if req.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
					}
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"message": "ok"}`)),
						Header:     make(http.Header),
					}
					return resp, nil
				},
			},
		},
	}

	ctx := context.Background()
	method := "GET"
	url := "/path"
	data := &bytes.Buffer{}

	result, err := c.Request(ctx, method, url, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedResult := `{"message": "ok"}`
	if string(result) != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, string(result))
	}
}

func TestClient_Request_NewRequestError_InvalidMethod(t *testing.T) {
	c := &Client{
		BaseUrl: "https://example.com",
	}

	ctx := context.Background()
	method := "INVALID METHOD"
	url := "/path"
	data := &bytes.Buffer{}

	_, err := c.Request(ctx, method, url, data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestClient_Request_DoError(t *testing.T) {
	c := &Client{
		BaseUrl: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("network error")
				},
			},
		},
	}

	ctx := context.Background()
	method := "GET"
	url := "/path"
	data := &bytes.Buffer{}

	_, err := c.Request(ctx, method, url, data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

type errorReadCloser struct {
	err error
}

func (e *errorReadCloser) Read(_ []byte) (n int, err error) {
	return 0, e.err
}

func (e *errorReadCloser) Close() error {
	return nil
}

func TestClient_Request_ReadAllError(t *testing.T) {
	c := &Client{
		BaseUrl: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: 200,
						Body: &errorReadCloser{
							err: errors.New("read error"),
						},
						Header: make(http.Header),
					}
					return resp, nil
				},
			},
		},
	}

	ctx := context.Background()
	method := "GET"
	url := "/path"
	data := &bytes.Buffer{}

	_, err := c.Request(ctx, method, url, data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

type errorCloseReadCloser struct {
	io.Reader
	err error
}

func (e *errorCloseReadCloser) Close() error {
	return e.err
}

func TestClient_Request_BodyCloseError(t *testing.T) {
	c := &Client{
		BaseUrl: "https://example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: 200,
						Body: &errorCloseReadCloser{
							Reader: strings.NewReader("response body"),
							err:    errors.New("close error"),
						},
						Header: make(http.Header),
					}
					return resp, nil
				},
			},
		},
	}

	ctx := context.Background()
	method := "GET"
	url := "/path"
	data := &bytes.Buffer{}

	result, err := c.Request(ctx, method, url, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedResult := "response body"
	if string(result) != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, string(result))
	}
}

func TestClient_Request_DefaultHTTPClient(t *testing.T) {
	c := &Client{
		BaseUrl:    "https://example.com",
		HTTPClient: nil,
	}

	originalTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"message": "ok"}`)),
				Header:     make(http.Header),
			}
			return resp, nil
		},
	}
	defer func() {
		http.DefaultTransport = originalTransport
	}()

	ctx := context.Background()
	method := "GET"
	url := "/path"
	data := &bytes.Buffer{}

	result, err := c.Request(ctx, method, url, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedResult := `{"message": "ok"}`
	if string(result) != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, string(result))
	}
}

func TestClient_resolveTrimmedUrl(t *testing.T) {
	tests := []struct {
		baseUrl string
		url     string
		want    string
	}{
		{"https://example.com/", "/path", "https://example.com/path"},
		{"https://example.com/", "path", "https://example.com/path"},
		{"https://example.com", "/path", "https://example.com/path"},
		{"https://example.com", "path", "https://example.com/path"},
		{"https://example.com/", "/", "https://example.com"},
		{"https://example.com", "", "https://example.com"},
		{"https://example.com/", "", "https://example.com"},
		{"https://example.com/", "//", "https://example.com"},
		{"https://example.com//", "/path", "https://example.com/path"},
	}

	for _, tt := range tests {
		c := &Client{
			BaseUrl: tt.baseUrl,
		}
		got := c.resolveTrimmedUrl(tt.url)
		if got != tt.want {
			t.Errorf("BaseUrl: %s, url: %s, expected: %s, got: %s", tt.baseUrl, tt.url, tt.want, got)
		}
	}
}
