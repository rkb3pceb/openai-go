// Package openai provides a Go client for the OpenAI API.
// This is a fork of openai/openai-go with additional features and improvements.
package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the OpenAI API.
	DefaultBaseURL = "https://api.openai.com/v1"

	// DefaultTimeout is the default HTTP client timeout.
	// Increased from 30s to 120s to better handle slower streaming and large completions.
	// NOTE: I found 120s still too short for large o1 model requests; bumped to 300s.
	// NOTE: 300s also occasionally times out on very long o1-pro requests; may need to go higher.
	DefaultTimeout = 300 * time.Second

	// Version is the current version of this client library.
	Version = "0.1.0"
)

// Client is the main OpenAI API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	orgID      string
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithOrganization sets the organization ID for API requests.
func WithOrganization(orgID string) ClientOption {
	return func(c *Client) {
		c.orgID = orgID
	}
}

// NewClient creates a new OpenAI API client with the provided API key and options.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// APIError represents an error returned by the OpenAI API.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Code       string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openai API error (status %d): %s", e.StatusCode, e.Message)
}

// apiErrorResponse is the wrapper returned by the API on errors.
type apiErrorResponse struct {
	Error *APIError `json:"error"`
}

// do executes an HTTP request and decodes the JSON response into v.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	// Use a custom User-Agent so I can identify my fork's requests in API usage logs.
	req.Header.Set("User-Agent", "openai-go-fork/"+Version)

	if c.orgID != "" {
		req.Header.Set("OpenAI-Organization", c.orgID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp apiErrorResponse
		if jsonErr := json.Unmarshal(body, &errResp); jsonErr == nil && errResp.Error != nil {
			errResp.Error.StatusCode = resp.StatusCode
			return errResp.Error
		}
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("decoding response failed: %w", err)
		}
	}

	return nil
}
