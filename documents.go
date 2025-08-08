package ragflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (c *Client) UploadDocument(ctx context.Context, datasetID, filePath string) (*Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("error copying file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents", datasetID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	var result ArrayResponse[Document]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no documents returned")
	}

	return &result.Data[0], nil
}

func (c *Client) UploadDocumentFromBytes(ctx context.Context, datasetID, filename string, data []byte) (*Document, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("error writing data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents", datasetID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	var result ArrayResponse[Document]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Println(result)
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no documents returned")
	}

	return &result.Data[0], nil
}

func (c *Client) GetDocument(ctx context.Context, datasetID, documentID string) (*Document, error) {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents/%s", datasetID, documentID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Document]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) ParseDocuments(ctx context.Context, datasetID string, documentIDs []string) error {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s/chunks", datasetID)
	httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, struct {
		IDs []string `json:"document_ids"`
	}{
		IDs: documentIDs,
		})
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

func (c *Client) DeleteDocuments(ctx context.Context, datasetID string, documentIDs []string) error {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents", datasetID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, struct {
		IDs []string `json:"ids"`
	}{
		IDs: documentIDs,
		})
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListDocumentsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Keywords string
	ID       string
}

func (c *Client) ListDocuments(ctx context.Context, datasetID string, opts *ListDocumentsOptions) (*DocumentsList, error) {
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
		if opts.Keywords != "" {
			params["keywords"] = opts.Keywords
		}
		if opts.ID != "" {
			params["id"] = opts.ID
		}
	}

	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents", datasetID)
	url := c.buildURL(endpoint, params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp DocumentsList
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DownloadDocument(ctx context.Context, datasetID, documentID string) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/datasets/%s/documents/%s/download", datasetID, documentID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, c.handleErrorResponse(resp.StatusCode, bodyBytes)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}

func (c *Client) GetChunk(ctx context.Context, chunkID string) (*Chunk, error) {
	endpoint := fmt.Sprintf("/api/v1/chunks/%s", chunkID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Chunk]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateChunk(ctx context.Context, chunkID string, req UpdateChunkRequest) (*Chunk, error) {
	endpoint := fmt.Sprintf("/api/v1/chunks/%s", chunkID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Chunk]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteChunk(ctx context.Context, chunkID string) error {
	endpoint := fmt.Sprintf("/api/v1/chunks/%s", chunkID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListChunksOptions struct {
	Page       int
	PageSize   int
	OrderBy    string
	Desc       bool
	Keywords   string
	ID         string
	DocumentID string
}

func (c *Client) ListChunks(ctx context.Context, datasetID string, opts *ListChunksOptions) (*ListResponse[Chunk], error) {
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
		if opts.Keywords != "" {
			params["keywords"] = opts.Keywords
		}
		if opts.ID != "" {
			params["id"] = opts.ID
		}
		if opts.DocumentID != "" {
			params["document_id"] = opts.DocumentID
		}
	}

	endpoint := fmt.Sprintf("/api/v1/datasets/%s/chunks", datasetID)
	url := c.buildURL(endpoint, params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Chunk]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
