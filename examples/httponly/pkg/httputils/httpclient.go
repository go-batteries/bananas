package httputils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Do(ctx context.Context, req *http.Request) (*HttpClientResponse, error)
	TlsEnabled() bool
}

// Transport
// Timeout
type HttpClient struct {
	// config HttpClientConfig
	client     *http.Client
	tlsEnabled bool
}

var DefaultHttpClientConfig = HttpClientConfig{
	Transport: http.DefaultTransport,
	Timeout:   10 * time.Second,
}

func NewHttpClient(config HttpClientConfig) *HttpClient {
	c := http.Client{
		Transport: config.Transport,
		Timeout:   config.Timeout,
	}

	return &HttpClient{client: &c, tlsEnabled: config.TlsEnabled}
}

type HttpClientConfig struct {
	Transport  http.RoundTripper
	Timeout    time.Duration
	TlsEnabled bool
}

type HttpClientResponse struct {
	Status  int
	Headers http.Header
	Body    io.Reader
}

func (slf *HttpClient) TlsEnabled() bool {
	return slf.tlsEnabled
}

// We are just going to go with the default transport for now.
// Adjust the MaxIdleConnnections later, we don't need 100 connections.
func (slf *HttpClient) Do(ctx context.Context, req *http.Request) (*HttpClientResponse, error) {
	resp, err := slf.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	var buf = bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpClientResponse{
		Status:  status,
		Headers: resp.Header.Clone(),
		Body:    buf,
	}, nil
}
