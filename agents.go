package ragflow

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (c *Client) CreateAgent(ctx context.Context, req CreateAgentRequest) (*Agent, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/v1/agents", req)
	if err != nil {
		return nil, err
	}

	var resp Response[Agent]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetAgent(ctx context.Context, agentID string) (*Agent, error) {
	endpoint := fmt.Sprintf("/api/v1/agents/%s", agentID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Agent]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateAgent(ctx context.Context, agentID string, req UpdateAgentRequest) (*Agent, error) {
	endpoint := fmt.Sprintf("/api/v1/agents/%s", agentID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Agent]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	endpoint := fmt.Sprintf("/api/v1/agents/%s", agentID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListAgentsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListAgents(ctx context.Context, opts *ListAgentsOptions) (*ListResponse[Agent], error) {
	params := make(map[string]string)
	
	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.OrderBy != "" {
			params["orderby"] = opts.OrderBy
		}
		if opts.Desc {
			params["desc"] = "true"
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.ID != "" {
			params["id"] = opts.ID
		}
	}

	url := c.buildURL("/api/v1/agents", params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Agent]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RunAgent(ctx context.Context, agentID string, message string, sessionID string) (*ChatCompletionResponse, error) {
	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
		ConversationID: sessionID,
	}

	endpoint := fmt.Sprintf("/api/v1/agents/%s/completions", agentID)
	httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp ChatCompletionResponse
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RunAgentStream(ctx context.Context, agentID string, message string, sessionID string) (<-chan ChatCompletionResponse, <-chan error) {
	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
		ConversationID: sessionID,
		Stream:         true,
	}

	respChan := make(chan ChatCompletionResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		endpoint := fmt.Sprintf("/api/v1/agents/%s/completions", agentID)
		httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, req)
		if err != nil {
			errChan <- err
			return
		}

		httpReq.Header.Set("Accept", "text/event-stream")
		httpReq.Header.Set("Cache-Control", "no-cache")

		resp, err := c.HTTPClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("error making streaming request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errChan <- c.handleErrorResponse(resp.StatusCode, bodyBytes)
			return
		}

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				errChan <- fmt.Errorf("error reading stream: %w", err)
				return
			}

			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if !bytes.HasPrefix(line, []byte("data: ")) {
				continue
			}

			data := bytes.TrimPrefix(line, []byte("data: "))
			if bytes.Equal(data, []byte("[DONE]")) {
				break
			}

			var streamResp ChatCompletionResponse
			if err := json.Unmarshal(data, &streamResp); err != nil {
				errChan <- fmt.Errorf("error unmarshaling stream data: %w", err)
				return
			}

			select {
			case respChan <- streamResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}
	}()

	return respChan, errChan
}