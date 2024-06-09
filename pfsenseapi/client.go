package pfsenseapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	Client interface {
		vlansClient
		interfacesClient
	}

	client struct {
		client  *http.Client
		options options
	}
)

func New(opts ...Options) (Client, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	if err := options.validate(); err != nil {
		return nil, err
	}

	newClient := &client{
		options: options,
		client: &http.Client{
			Timeout: options.timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: options.skipTLS},
			},
		},
	}
	newClient.client.CloseIdleConnections()

	return newClient, nil
}

func (c *client) do(ctx context.Context, method, endpoint string, queryMap map[string]string, body []byte) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.options.host, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range queryMap {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.options.username, c.options.password)

	return c.client.Do(req)
}

func (c *client) get(ctx context.Context, endpoint string, queryMap map[string]string) ([]byte, error) {
	res, err := c.do(ctx, http.MethodGet, endpoint, queryMap, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	respbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err, ok := responseCodeErrorMap[res.StatusCode]
		if !ok {
			err = fmt.Errorf("non 2xx response code received: %d", res.StatusCode)
		}

		resp := new(apiResponse)
		if jsonerr := json.Unmarshal(respbody, resp); jsonerr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, resp.Message)
	}

	return respbody, nil
}

func (c *client) post(ctx context.Context, endpoint string, queryMap map[string]string, body []byte) ([]byte, error) {
	res, err := c.do(ctx, http.MethodPost, endpoint, queryMap, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	respbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err, ok := responseCodeErrorMap[res.StatusCode]
		if !ok {
			err = fmt.Errorf("non 2xx response code received: %d", res.StatusCode)
		}

		resp := new(apiResponse)
		if jsonerr := json.Unmarshal(respbody, resp); jsonerr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, resp.Message)
	}

	return respbody, nil
}

// patch makes a patch request to the given endpoint with the given queryMap and body.
func (c *client) patch(ctx context.Context, endpoint string, queryMap map[string]string, body []byte) ([]byte, error) {
	res, err := c.do(ctx, http.MethodPatch, endpoint, queryMap, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close
	}()

	respbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err, ok := responseCodeErrorMap[res.StatusCode]
		if !ok {
			err = fmt.Errorf("non 2xx response code received: %d", res.StatusCode)
		}

		resp := new(apiResponse)
		if jsonerr := json.Unmarshal(respbody, resp); jsonerr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, resp.Message)
	}

	return respbody, nil
}

func (c *client) put(ctx context.Context, endpoint string, queryMap map[string]string, body []byte) ([]byte, error) {
	res, err := c.do(ctx, http.MethodPut, endpoint, queryMap, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	respbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err, ok := responseCodeErrorMap[res.StatusCode]
		if !ok {
			err = fmt.Errorf("non 2xx response code received: %d", res.StatusCode)
		}

		resp := new(apiResponse)
		if jsonerr := json.Unmarshal(respbody, resp); jsonerr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, resp.Message)
	}

	return respbody, nil
}

func (c *client) delete(ctx context.Context, endpoint string, queryMap map[string]string) ([]byte, error) {
	res, err := c.do(ctx, http.MethodDelete, endpoint, queryMap, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	respbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err, ok := responseCodeErrorMap[res.StatusCode]
		if !ok {
			err = fmt.Errorf("non 2xx response code received: %d", res.StatusCode)
		}

		resp := new(apiResponse)
		if jsonerr := json.Unmarshal(respbody, resp); jsonerr != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, resp.Message)
	}

	return respbody, nil
}

type apiResponse struct {
	Status     string `json:"status"`
	Code       int    `json:"code"`
	ResponseId string `json:"response_id"`
	Message    string `json:"message"`
}

var (
	// ErrBadRequest represents a HTTP 400 error
	ErrBadRequest = fmt.Errorf("HTTP 400: Bad Request")

	// ErrUnauthorized represents a HTTP 401 error
	ErrUnauthorized = fmt.Errorf("HTTP 401: Unauthorized")

	// ErrForbidden represents a HTTP 403 error
	ErrForbidden = fmt.Errorf("HTTP 403: Forbidden")

	// ErrNotFound represents a HTTP 404 error
	ErrNotFound = fmt.Errorf("HTTP 404: Not Found")

	// ErrMethodNotAllowed represents a HTTP 405 error
	ErrMethodNotAllowed = fmt.Errorf("HTTP 405: Method Not Allowed")

	// ErrNotAcceptable represents a HTTP 406 error
	ErrNotAcceptable = fmt.Errorf("HTTP 406: Not Acceptable")

	// ErrConflict represents a HTTP 409 error
	ErrConflict = fmt.Errorf("HTTP 409: Conflict")

	// ErrUnsupportedMediaType represents a HTTP 415 error
	ErrUnsupportedMediaType = fmt.Errorf("HTTP 415: Unsupported Media Type")

	// ErrUnprocessableEntity represents a HTTP 422 error
	ErrUnprocessableEntity = fmt.Errorf("HTTP 422: Unprocessable Entity")

	// ErrFailedDependency represents a HTTP 424 error
	ErrFailedDependency = fmt.Errorf("HTTP 424: Failed Dependency")

	// ErrInternalServerError represents a HTTP 500 error
	ErrInternalServerError = fmt.Errorf("HTTP 500: Internal Server Error")

	// ErrServiceUnavailable represents a HTTP 503 error
	ErrServiceUnavailable = fmt.Errorf("HTTP 503: Service Unavailable")

	responseCodeErrorMap = map[int]error{
		http.StatusBadRequest:           ErrBadRequest,
		http.StatusUnauthorized:         ErrUnauthorized,
		http.StatusForbidden:            ErrForbidden,
		http.StatusNotFound:             ErrNotFound,
		http.StatusMethodNotAllowed:     ErrMethodNotAllowed,
		http.StatusNotAcceptable:        ErrNotAcceptable,
		http.StatusConflict:             ErrConflict,
		http.StatusUnsupportedMediaType: ErrUnsupportedMediaType,
		http.StatusUnprocessableEntity:  ErrUnprocessableEntity,
		http.StatusFailedDependency:     ErrFailedDependency,
		http.StatusInternalServerError:  ErrInternalServerError,
		http.StatusServiceUnavailable:   ErrServiceUnavailable,
	}
)
