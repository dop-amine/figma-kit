//go:build mcp_go_client_oauth

package mcpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const figmaMCPEndpoint = "https://mcp.figma.com/mcp"

// Session wraps an MCP client session connected to the Figma MCP server.
type Session struct {
	cs *mcp.ClientSession
}

// Connect establishes a connection to the Figma MCP server with OAuth.
func Connect(ctx context.Context) (*Session, error) {
	authChan := make(chan *auth.AuthorizationResult, 1)
	errChan := make(chan error, 1)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("finding free port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://localhost:%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authChan <- &auth.AuthorizationResult{
			Code:  r.URL.Query().Get("code"),
			State: r.URL.Query().Get("state"),
		}
		fmt.Fprint(w, "<html><body><h2>Authentication successful!</h2><p>You can close this window and return to the terminal.</p></body></html>")
	})
	srv := &http.Server{Handler: mux}
	go func() {
		if sErr := srv.Serve(listener); sErr != nil && sErr != http.ErrServerClosed {
			errChan <- sErr
		}
	}()
	defer srv.Close()

	codeFetcher := func(_ context.Context, args *auth.AuthorizationArgs) (*auth.AuthorizationResult, error) {
		fmt.Printf("Opening browser for Figma authentication...\n")
		if bErr := openSystemBrowser(args.URL); bErr != nil {
			fmt.Printf("Could not open browser automatically.\nPlease visit: %s\n", args.URL)
		}
		select {
		case res := <-authChan:
			return res, nil
		case e := <-errChan:
			return nil, e
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	authHandler, err := auth.NewAuthorizationCodeHandler(&auth.AuthorizationCodeHandlerConfig{
		RedirectURL:              redirectURL,
		AuthorizationCodeFetcher: codeFetcher,
	})
	if err != nil {
		return nil, fmt.Errorf("creating auth handler: %w", err)
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint:     figmaMCPEndpoint,
		OAuthHandler: authHandler,
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "figma-kit",
		Version: "1.0.0",
	}, nil)

	cs, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to Figma MCP: %w", err)
	}

	return &Session{cs: cs}, nil
}

// Close terminates the MCP session.
func (s *Session) Close() error {
	if s.cs != nil {
		return s.cs.Close()
	}
	return nil
}

// UseFigmaResult is the result of a use_figma call.
type UseFigmaResult struct {
	Content []string
	IsError bool
}

// CallUseFigma executes JavaScript in a Figma file via the use_figma tool.
func (s *Session) CallUseFigma(ctx context.Context, fileKey, code, description string) (*UseFigmaResult, error) {
	params := &mcp.CallToolParams{
		Name: "use_figma",
		Arguments: map[string]any{
			"fileKey":     fileKey,
			"code":        code,
			"description": description,
		},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("use_figma call failed: %w", err)
	}
	result := &UseFigmaResult{IsError: res.IsError}
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			result.Content = append(result.Content, tc.Text)
		}
	}
	return result, nil
}

// WhoamiResult holds the whoami response.
type WhoamiResult struct {
	Raw string
}

// CallWhoami returns information about the authenticated Figma user.
func (s *Session) CallWhoami(ctx context.Context) (*WhoamiResult, error) {
	params := &mcp.CallToolParams{
		Name:      "whoami",
		Arguments: map[string]any{},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return nil, err
	}
	var texts []string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			texts = append(texts, tc.Text)
		}
	}
	return &WhoamiResult{Raw: joinTexts(texts)}, nil
}

// CallScreenshot captures a screenshot of a Figma node.
func (s *Session) CallScreenshot(ctx context.Context, fileKey, nodeID string) (string, error) {
	params := &mcp.CallToolParams{
		Name: "get_screenshot",
		Arguments: map[string]any{
			"fileKey": fileKey,
			"nodeId":  nodeID,
		},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return "", err
	}
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			return tc.Text, nil
		}
	}
	return "", fmt.Errorf("no text content in screenshot response")
}

// FileResult holds create_new_file response data.
type FileResult struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
	Raw     string
}

// CallCreateFile creates a new Figma file.
func (s *Session) CallCreateFile(ctx context.Context, name, planKey string) (*FileResult, error) {
	params := &mcp.CallToolParams{
		Name: "create_new_file",
		Arguments: map[string]any{
			"fileName":   name,
			"planKey":    planKey,
			"editorType": "design",
		},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return nil, err
	}
	result := &FileResult{}
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			result.Raw = tc.Text
			_ = json.Unmarshal([]byte(tc.Text), result)
		}
	}
	return result, nil
}

// SearchResult holds search_design_system response.
type SearchResult struct {
	Raw string
}

// CallSearchDS searches the design system for components, variables, and styles.
func (s *Session) CallSearchDS(ctx context.Context, fileKey, query string) (*SearchResult, error) {
	params := &mcp.CallToolParams{
		Name: "search_design_system",
		Arguments: map[string]any{
			"fileKey": fileKey,
			"query":   query,
		},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return nil, err
	}
	var texts []string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			texts = append(texts, tc.Text)
		}
	}
	return &SearchResult{Raw: joinTexts(texts)}, nil
}

// CallGetMetadata retrieves metadata (XML structure) for a node.
func (s *Session) CallGetMetadata(ctx context.Context, fileKey, nodeID string) (string, error) {
	params := &mcp.CallToolParams{
		Name: "get_metadata",
		Arguments: map[string]any{
			"fileKey": fileKey,
			"nodeId":  nodeID,
		},
	}
	res, err := s.cs.CallTool(ctx, params)
	if err != nil {
		return "", err
	}
	var texts []string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			texts = append(texts, tc.Text)
		}
	}
	return joinTexts(texts), nil
}

func joinTexts(texts []string) string {
	result := ""
	for i, t := range texts {
		if i > 0 {
			result += "\n"
		}
		result += t
	}
	return result
}

func openSystemBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	}
	return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}
