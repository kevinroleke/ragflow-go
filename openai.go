package ragflow

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

		endpoint := fmt.Sprintf("/api/v1/chats_openai/%s/chat/completions", req.ConversationID)

		req.Stream = true

		httpReq, err := c.newRequest(ctx, "POST", endpoint, req)
		if err != nil {
			errChan <- err
			log.Println("46new")
			return
		}

		resp, err := c.HTTPClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("error making request: %w", err)
			log.Println("51req")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errChan <- c.handleErrorResponse(resp.StatusCode, bodyBytes)
			log.Println("59over400")
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				log.Println("72[DONE]")
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
				log.Println("ctx done errchan")
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
		}
	}()

	return respChan, errChan
}
