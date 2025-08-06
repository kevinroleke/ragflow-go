package ragflow

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) CreateDataset(ctx context.Context, req CreateDatasetRequest) (*Dataset, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/v1/datasets", req)
	if err != nil {
		return nil, err
	}

	var resp Response[Dataset]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetDataset(ctx context.Context, datasetID string) (*Dataset, error) {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Dataset]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateDataset(ctx context.Context, datasetID string, req UpdateDatasetRequest) (*Dataset, error) {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Dataset]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteDataset(ctx context.Context, datasetID string) error {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListDatasetsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListDatasets(ctx context.Context, opts *ListDatasetsOptions) (*ListResponse[Dataset], error) {
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

	url := c.buildURL("/api/v1/datasets", params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Dataset]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
