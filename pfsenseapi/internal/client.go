package internal

import (
	"io"
	"net/http"
	"net/url"
)

var _ Client = &http.Client{}
var _ Client = (*client)(nil)

type Client interface {
	CloseIdleConnections()
	Do(req *Request) (*http.Response, error)
	Head(url string) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type client struct {
	http.Client
	pipeline Pipeline
}
