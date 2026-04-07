package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const baseURL = "https://api.figma.com"

// Client wraps the Figma REST API. Created via NewClient(), which returns nil
// when no PAT is available — callers should check for nil before using.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a REST API client if a PAT is available.
// Returns nil if no PAT is set (FIGMA_TOKEN, FIGMA_PAT).
func NewClient() *Client {
	pat := os.Getenv("FIGMA_TOKEN")
	if pat == "" {
		pat = os.Getenv("FIGMA_PAT")
	}
	if pat == "" {
		pat = os.Getenv("FIGMA_PERSONAL_ACCESS_TOKEN")
	}
	if pat == "" {
		return nil
	}
	return &Client{
		token:      pat,
		httpClient: &http.Client{},
	}
}

// Available returns true if the REST API client is usable.
func (c *Client) Available() bool {
	return c != nil && c.token != ""
}

// FileMeta holds basic file metadata.
type FileMeta struct {
	Name         string `json:"name"`
	LastModified string `json:"lastModified"`
	Version      string `json:"version"`
	Role         string `json:"role"`
}

// GetFileMeta retrieves metadata for a Figma file.
func (c *Client) GetFileMeta(fileKey string) (*FileMeta, error) {
	body, err := c.get(fmt.Sprintf("/v1/files/%s/meta", fileKey))
	if err != nil {
		return nil, err
	}
	var result FileMeta
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing file meta: %w", err)
	}
	return &result, nil
}

// FileJSON is the top-level structure returned by GET /v1/files/:key.
type FileJSON struct {
	Name     string          `json:"name"`
	Document json.RawMessage `json:"document"`
	Version  string          `json:"version"`
}

// GetFile retrieves the full file JSON.
func (c *Client) GetFile(fileKey string) (*FileJSON, error) {
	body, err := c.get(fmt.Sprintf("/v1/files/%s", fileKey))
	if err != nil {
		return nil, err
	}
	var result FileJSON
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}
	return &result, nil
}

// NodesJSON is the response from GET /v1/files/:key/nodes.
type NodesJSON struct {
	Nodes map[string]json.RawMessage `json:"nodes"`
}

// GetNodes retrieves specific node subtrees.
func (c *Client) GetNodes(fileKey string, ids []string) (*NodesJSON, error) {
	idsParam := strings.Join(ids, ",")
	body, err := c.get(fmt.Sprintf("/v1/files/%s/nodes?ids=%s", fileKey, idsParam))
	if err != nil {
		return nil, err
	}
	var result NodesJSON
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing nodes: %w", err)
	}
	return &result, nil
}

// ExportImages renders nodes as images. Returns a map of node ID → image URL.
func (c *Client) ExportImages(fileKey string, ids []string, format string, scale float64) (map[string]string, error) {
	idsParam := strings.Join(ids, ",")
	if format == "" {
		format = "png"
	}
	if scale <= 0 {
		scale = 1
	}
	body, err := c.get(fmt.Sprintf("/v1/images/%s?ids=%s&format=%s&scale=%g", fileKey, idsParam, format, scale))
	if err != nil {
		return nil, err
	}
	var result struct {
		Images map[string]string `json:"images"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing images: %w", err)
	}
	return result.Images, nil
}

// FileImages retrieves all image fills in a file. Returns a map of image hash → URL.
func (c *Client) FileImages(fileKey string) (map[string]string, error) {
	body, err := c.get(fmt.Sprintf("/v1/files/%s/images", fileKey))
	if err != nil {
		return nil, err
	}
	var result struct {
		Meta struct {
			Images map[string]string `json:"images"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing file images: %w", err)
	}
	return result.Meta.Images, nil
}

func (c *Client) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Figma-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
