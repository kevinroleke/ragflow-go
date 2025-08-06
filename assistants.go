package ragflow

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) CreateAssistant(ctx context.Context, req CreateAssistantRequest) (*Assistant, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/v1/chats", req)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetAssistant(ctx context.Context, assistantID string) (*Assistant, error) {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateAssistant(ctx context.Context, assistantID string, req UpdateAssistantRequest) (*Assistant, error) {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListAssistantsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListAssistants(ctx context.Context, opts *ListAssistantsOptions) (*ListResponse[Assistant], error) {
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

	url := c.buildURL("/api/v1/chat/assistants", params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateSession(ctx context.Context, assistantID string, req CreateSessionRequest) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetSession(ctx context.Context, assistantID, sessionID string) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateSession(ctx context.Context, assistantID, sessionID string, req UpdateSessionRequest) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteSession(ctx context.Context, assistantID, sessionID string) error {
	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListSessionsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListSessions(ctx context.Context, assistantID string, opts *ListSessionsOptions) (*ListResponse[Session], error) {
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

	endpoint := fmt.Sprintf("/api/v1/chat/assistants/%s/sessions", assistantID)
	url := c.buildURL(endpoint, params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
