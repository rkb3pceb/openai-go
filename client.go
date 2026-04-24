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
	DefaultTimeout = 30 * time.Second

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
	req.Header.Set("User-Agent", "openai-go/"+Version)

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
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr apiErrorResponse
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.Error != nil {
			apiErr.Error.StatusCode = resp.StatusCode
			return apiErr.Error
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}
