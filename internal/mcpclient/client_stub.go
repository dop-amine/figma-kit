//go:build !mcp_go_client_oauth

package mcpclient

import (
	"context"
	"fmt"
)

// Session wraps an MCP client session connected to the Figma MCP server.
type Session struct{}

// Connect establishes a connection to the Figma MCP server.
// Without the mcp_go_client_oauth build tag, this always returns an error
// directing the user to build with the tag or use the AI agent workflow.
func Connect(_ context.Context) (*Session, error) {
	return nil, fmt.Errorf(
		"direct MCP execution not available in this build\n\n" +
			"To enable: rebuild with 'go build -tags mcp_go_client_oauth'\n" +
			"Or use the AI agent workflow: commands output JS for use_figma")
}

// Close terminates the MCP session.
func (s *Session) Close() error { return nil }

// UseFigmaResult is the result of a use_figma call.
type UseFigmaResult struct {
	Content []string
	IsError bool
}

// CallUseFigma executes JavaScript in a Figma file.
func (s *Session) CallUseFigma(_ context.Context, _, _, _ string) (*UseFigmaResult, error) {
	return nil, fmt.Errorf("MCP not available in this build")
}

// WhoamiResult holds the whoami response.
type WhoamiResult struct {
	Raw string
}

// CallWhoami returns information about the authenticated Figma user.
func (s *Session) CallWhoami(_ context.Context) (*WhoamiResult, error) {
	return nil, fmt.Errorf("MCP not available in this build")
}

// CallScreenshot captures a screenshot of a Figma node.
func (s *Session) CallScreenshot(_ context.Context, _, _ string) (string, error) {
	return "", fmt.Errorf("MCP not available in this build")
}

// FileResult holds create_new_file response data.
type FileResult struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
	Raw     string
}

// CallCreateFile creates a new Figma file.
func (s *Session) CallCreateFile(_ context.Context, _, _ string) (*FileResult, error) {
	return nil, fmt.Errorf("MCP not available in this build")
}

// SearchResult holds search_design_system response.
type SearchResult struct {
	Raw string
}

// CallSearchDS searches the design system.
func (s *Session) CallSearchDS(_ context.Context, _, _ string) (*SearchResult, error) {
	return nil, fmt.Errorf("MCP not available in this build")
}

// CallGetMetadata retrieves metadata for a node.
func (s *Session) CallGetMetadata(_ context.Context, _, _ string) (string, error) {
	return "", fmt.Errorf("MCP not available in this build")
}
