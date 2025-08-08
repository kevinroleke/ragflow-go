package ragflow

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
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
	Username   string
	Password   string
	SessionCookie string
	SessionAuth string
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

func WithUserPass(username string, password string) ClientOption {
	return func(c *Client) {
		c.Username = username
		c.Password = password
		auth, cookie, err := c.Login(context.Background())
		if err != nil {
			panic(err)
		}
		c.SessionCookie = cookie
		c.SessionAuth = auth
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

func rsaPsw(password string) (string, error) {
    // The same public key from the JavaScript code
    pubKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArq9XTUSeYr2+N1h3Afl/z8Dse/2yD0ZGrKwx+EEEcdsBLca9Ynmx3nIB5obmLlSfmskLpBo0UACBmB5rEjBp2Q2f3AG3Hjd4B+gNCG6BDaawuDlgANIhGnaTLrIqWrrcm4EMzJOnAOI1fgzJRsOOUEfaS318Eq9OVO3apEyCCt0lOQK6PuksduOjVxtltDav+guVAA068NrPYmRNabVKRNLJpL8w4D44sfth5RvZ3q9t+6RTArpEtc5sh5ChzvqPOzKGMXW83C95TxmXqpbK6olN4RevSfVjEAgCydH6HN6OhtOQEcnrU97r9H0iZOWwbw3pVrZiUkuRD1R56Wzs2wIDAQAB
-----END PUBLIC KEY-----`

    // Parse PEM block
    block, _ := pem.Decode([]byte(pubKeyPEM))
    if block == nil {
        return "", fmt.Errorf("failed to parse PEM block")
    }

    // Parse public key
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return "", fmt.Errorf("failed to parse public key: %v", err)
    }

    // Type assert to RSA public key
    rsaPub, ok := pub.(*rsa.PublicKey)
    if !ok {
        return "", fmt.Errorf("not an RSA public key")
    }

    // Base64 encode the password (matching Base64.encode(password) in JS)
    encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

    // RSA encrypt using PKCS1v15 (same as JSEncrypt default)
    encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(encodedPassword))
    if err != nil {
        return "", fmt.Errorf("encryption failed: %v", err)
    }

    // Return base64 encoded result (matching JSEncrypt output)
    return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (c *Client) Login(ctx context.Context) (string, string, error) {
	url := c.BaseURL + "/v1/user/login"
	rsaPassword, err := rsaPsw(c.Password)
	if err != nil {
		panic(err)
	}

	body := struct{
		Username string `json:"email"`
		Password string `json:"password"`
	}{
		Username: c.Username,
		Password: rsaPassword,
	}

	var buf io.Reader
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling request body: %w", err)
	}
	buf = bytes.NewBuffer(jsonData)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", "", c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}
	// Check for API-level errors in the response
	if err := c.checkAPIResponse(bodyBytes); err != nil {
		return "", "", err
	}

	auth := resp.Header.Get("Authorization")
	cookie := resp.Header.Get("Set-Cookie")
	if auth == "" {
		return "", "", fmt.Errorf("Fail to login")
	}
	if cookie == "" {
		return "", "", fmt.Errorf("Fail to login")
	}

	return auth, strings.Split(strings.Split(cookie, "session=")[1], ";")[0], nil
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
	u := endpoint
	if len(params) == 0 {
		return u
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return u + "?" + values.Encode()
}
