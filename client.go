package ragflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
	DefaultBaseURL = "http://127.0.0.1"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{}
		}
		c.HTTPClient.Timeout = timeout
	}
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	url := c.BaseURL + endpoint

	var buf io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	if v != nil {
		if err := json.Unmarshal(bodyBytes, v); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
		
		// Check for API-level errors in the response
		if err := c.checkAPIResponse(bodyBytes); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			Code:       statusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", statusCode, string(body)),
			StatusCode: statusCode,
		}
	}

	return &APIError{
		Code:       errResp.Code,
		Message:    errResp.Message,
		StatusCode: statusCode,
	}
}

func (c *Client) checkAPIResponse(body []byte) error {
	var baseResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	
	if err := json.Unmarshal(body, &baseResp); err != nil {
		return nil // If we can't parse it, assume it's not an error response
	}
	
	// Check if the API returned an error code
	if baseResp.Code != 0 && baseResp.Code != 200 {
		return &APIError{
			Code:       baseResp.Code,
			Message:    baseResp.Message,
			StatusCode: 200, // HTTP was OK but API returned error
		}
	}
	
	return nil
}

func (c *Client) buildURL(endpoint string, params map[string]string) string {
	u := c.BaseURL + endpoint
	if len(params) == 0 {
		return u
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return u + "?" + values.Encode()
}