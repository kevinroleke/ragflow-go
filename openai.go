package ragflow

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/chats_openai/%s/chat/completions", req.Model)

	req.Stream = false

	httpReq, err := c.newRequest(ctx, "POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response ChatCompletionResponse
	if err := c.do(httpReq, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan ChatCompletionResponse, <-chan error) {
	respChan := make(chan ChatCompletionResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		endpoint := fmt.Sprintf("/api/v1/chats_openai/%s/chat/completions", req.Model)

		req.Stream = true

		httpReq, err := c.newRequest(ctx, "POST", endpoint, req)
		if err != nil {
			errChan <- err
			return
		}

		resp, err := c.HTTPClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("error making request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errChan <- c.handleErrorResponse(resp.StatusCode, bodyBytes)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" {
				continue
			}

			var data string
			if strings.HasPrefix(line, "data:") {
				data = strings.TrimPrefix(line, "data:")
			} else {
				// Handle non-SSE formatted responses (plain JSON)
				data = line
			}

			if data == "[DONE]" {
				return
			}

			// Check if this is an error response (non-streaming)
			var errorCheck struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal([]byte(data), &errorCheck); err == nil && errorCheck.Code != 0 {
				errChan <- fmt.Errorf("API error %d: %s", errorCheck.Code, errorCheck.Message)
				return
			}

			var streamResp ChatCompletionResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			select {
			case respChan <- streamResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
		}
	}()

	return respChan, errChan
}
