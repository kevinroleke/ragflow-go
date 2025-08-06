package ragflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LLMModel struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	UsedToken int    `json:"used_token"`
}

// Provider grouping of LLMs with associated tags
type LLMProvider struct {
	LLMs []LLMModel `json:"llm"`
	Tags string     `json:"tags"`
}

// Response structure for /v1/llm/my_llms endpoint
type MyLLMsResponse map[string]LLMProvider

type Factory struct {
	CreateDate  string   `json:"create_date"`
	CreateTime  int64    `json:"create_time"`
	Logo        string   `json:"logo"`
	ModelTypes  []string `json:"model_types"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Tags        string   `json:"tags"`
	UpdateDate  string   `json:"update_date"`
	UpdateTime  int64    `json:"update_time"`
}

func (c *Client) newUserRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
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

	if c.SessionAuth != "" {
		req.Header.Set("Authorization", c.SessionAuth)
	}
	if c.SessionCookie != "" {
		req.Header.Set("Cookie", "session="+c.SessionCookie)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) GetMyLLMs(ctx context.Context) (MyLLMsResponse, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/v1/llm/my_llms", nil)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()
	bytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var response Response[MyLLMsResponse]
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}

func (c *Client) GetFactories(ctx context.Context) ([]Factory, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/v1/llm/factories", nil)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()
	bytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var response Response[[]Factory]
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}
