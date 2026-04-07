package mcpclient

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const figmaMCPEndpoint = "https://mcp.figma.com/mcp"

// Session wraps an MCP client session connected to the Figma MCP server.
type Session struct {
	cs *mcp.ClientSession
}

// Connect establishes a connection to the Figma MCP server.
// It loads a cached token or runs the full OAuth flow (PAT → register → PKCE → token).
func Connect(ctx context.Context) (*Session, error) {
	token, err := loadOrFetchToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint:   figmaMCPEndpoint,
		HTTPClient: tokenHTTPClient(token.AccessToken),
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

// tokenHTTPClient returns an http.Client that injects Bearer token on every request.
func tokenHTTPClient(accessToken string) *http.Client {
	return &http.Client{
		Transport: &bearerTransport{
			token: accessToken,
			base:  http.DefaultTransport,
		},
	}
}

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(r)
}

// loadOrFetchToken loads a token using a priority chain:
// 1. Cached token from ~/.config/figma-kit/token.json
// 2. FIGMA_ACCESS_TOKEN env var (direct OAuth token, no registration needed)
// 3. Full OAuth flow using PAT for client registration
func loadOrFetchToken(ctx context.Context) (*TokenData, error) {
	// Tier 1: cached token
	cached, err := LoadToken()
	if err == nil && cached.IsValid() {
		return cached, nil
	}

	// Tier 2: direct OAuth token from environment
	if accessToken := os.Getenv("FIGMA_ACCESS_TOKEN"); accessToken != "" {
		td := &TokenData{AccessToken: accessToken}
		return td, nil
	}

	// Tier 3: full OAuth flow (PAT → register → PKCE → token)
	return runOAuthFlow(ctx)
}

// oauthServerMeta holds discovered OAuth server metadata.
type oauthServerMeta struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	RegistrationEndpoint  string `json:"registration_endpoint"`
}

// registeredClient holds the result of dynamic client registration.
type registeredClient struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// runOAuthFlow performs the full OAuth 2.0 flow:
// 1. Get a Figma PAT (env var or prompt)
// 2. Discover OAuth endpoints
// 3. Register a dynamic client (using PAT for auth)
// 4. PKCE authorization code flow
// 5. Exchange code for access + refresh tokens
func runOAuthFlow(ctx context.Context) (*TokenData, error) {
	// Step 1: Get PAT
	pat := getFigmaPAT()
	if pat == "" {
		return nil, fmt.Errorf("Figma Personal Access Token required for first-time setup\n\n" +
			"Get one at: Figma → Settings → Security → Personal access tokens\n" +
			"Then either:\n" +
			"  export FIGMA_PAT=your_token\n" +
			"  figma-kit auth login\n" +
			"Or run 'figma-kit auth login' and paste when prompted")
	}

	// Step 2: Discover OAuth metadata
	fmt.Println("Discovering Figma OAuth endpoints...")
	meta, err := discoverOAuthMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("OAuth discovery: %w", err)
	}

	// Step 3: Register dynamic client
	fmt.Println("Registering OAuth client...")
	client, err := registerClient(ctx, meta, pat)
	if err != nil {
		return nil, fmt.Errorf("client registration: %w", err)
	}

	// Step 4: Start local redirect server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("starting local server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	codeCh := make(chan authCallback, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		codeCh <- authCallback{
			code:  r.URL.Query().Get("code"),
			state: r.URL.Query().Get("state"),
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body style="font-family:system-ui;padding:3em;text-align:center">
			<h2 style="color:#33CC80">✓ Authenticated!</h2>
			<p>You can close this tab and return to the terminal.</p></body></html>`)
	})
	srv := &http.Server{Handler: mux}
	go srv.Serve(listener)
	defer srv.Close()

	// Step 5: PKCE challenge
	verifier := generateVerifier()
	challenge := generateS256Challenge(verifier)
	state := generateState()

	authURL := buildAuthURL(meta.AuthorizationEndpoint, client.ClientID, redirectURL, state, challenge)

	fmt.Println("Opening browser for Figma authorization...")
	if bErr := openSystemBrowser(authURL); bErr != nil {
		fmt.Printf("\nCould not open browser automatically.\nPlease visit:\n  %s\n\n", authURL)
	}

	// Step 6: Wait for callback
	var cb authCallback
	select {
	case cb = <-codeCh:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timed out after 5 minutes")
	}

	if cb.state != state {
		return nil, fmt.Errorf("OAuth state mismatch — possible CSRF attack")
	}
	if cb.code == "" {
		return nil, fmt.Errorf("no authorization code received")
	}

	// Step 7: Exchange code for token
	fmt.Println("Exchanging authorization code for token...")
	tokenData, err := exchangeCode(ctx, meta.TokenEndpoint, client.ClientID, client.ClientSecret, cb.code, redirectURL, verifier)
	if err != nil {
		return nil, err
	}

	if err := SaveToken(tokenData); err != nil {
		fmt.Printf("Warning: could not cache token: %v\n", err)
	}

	fmt.Println("✓ Authenticated with Figma successfully!")
	return tokenData, nil
}

type authCallback struct {
	code  string
	state string
}

// getFigmaPAT reads the Figma Personal Access Token from env or prompts.
// Checks FIGMA_TOKEN first (common env var name), then FIGMA_PAT, then
// FIGMA_PERSONAL_ACCESS_TOKEN, then interactive prompt.
func getFigmaPAT() string {
	if pat := os.Getenv("FIGMA_TOKEN"); pat != "" {
		return pat
	}
	if pat := os.Getenv("FIGMA_PAT"); pat != "" {
		return pat
	}
	if pat := os.Getenv("FIGMA_PERSONAL_ACCESS_TOKEN"); pat != "" {
		return pat
	}
	// Check if we have stdin (interactive terminal)
	fi, err := os.Stdin.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice) == 0 {
		return ""
	}
	fmt.Println()
	fmt.Println("First-time setup: Figma Personal Access Token required.")
	fmt.Println("Get one at: Figma → Settings → Security → Personal access tokens")
	fmt.Println()
	fmt.Print("Paste your Figma PAT: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// discoverOAuthMetadata fetches the OAuth server metadata from Figma.
func discoverOAuthMetadata(ctx context.Context) (*oauthServerMeta, error) {
	// First discover the protected resource metadata
	prmURL := "https://mcp.figma.com/.well-known/oauth-protected-resource/mcp"
	var authServerBase string

	req, err := http.NewRequestWithContext(ctx, "GET", prmURL, nil)
	if err == nil {
		resp, respErr := http.DefaultClient.Do(req)
		if respErr == nil && resp.StatusCode == 200 {
			var prm struct {
				AuthorizationServers []string `json:"authorization_servers"`
			}
			json.NewDecoder(resp.Body).Decode(&prm)
			resp.Body.Close()
			if len(prm.AuthorizationServers) > 0 {
				authServerBase = prm.AuthorizationServers[0]
			}
		} else if resp != nil {
			resp.Body.Close()
		}
	}

	if authServerBase == "" {
		authServerBase = "https://api.figma.com"
	}

	// Fetch authorization server metadata
	asmURL := authServerBase + "/.well-known/oauth-authorization-server"
	req, err = http.NewRequestWithContext(ctx, "GET", asmURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching auth server metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("auth server metadata returned HTTP %d", resp.StatusCode)
	}

	var meta oauthServerMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("parsing auth server metadata: %w", err)
	}

	return &meta, nil
}

// registerClient performs Figma's dynamic client registration.
// Figma requires a PAT in the X-Figma-Token header for registration.
func registerClient(ctx context.Context, meta *oauthServerMeta, pat string) (*registeredClient, error) {
	endpoint := meta.RegistrationEndpoint
	if endpoint == "" {
		endpoint = "https://api.figma.com/v1/oauth/mcp/register"
	}

	regBody := map[string]any{
		"client_name":                "figma-kit",
		"redirect_uris":             []string{"http://127.0.0.1"},
		"grant_types":               []string{"authorization_code", "refresh_token"},
		"response_types":            []string{"code"},
		"token_endpoint_auth_method": "none",
		"scope":                     "mcp:connect",
	}
	body, _ := json.Marshal(regBody)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Figma-Token", pat)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("registration forbidden — check that your Figma PAT is valid")
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result registeredClient
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parsing registration response: %w", err)
	}
	if result.ClientID == "" {
		return nil, fmt.Errorf("registration succeeded but no client_id returned")
	}

	return &result, nil
}

func buildAuthURL(endpoint, clientID, redirectURL, state, challenge string) string {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", clientID)
	v.Set("redirect_uri", redirectURL)
	v.Set("state", state)
	v.Set("code_challenge", challenge)
	v.Set("code_challenge_method", "S256")
	v.Set("scope", "mcp:connect")
	return endpoint + "?" + v.Encode()
}

// exchangeCode exchanges an authorization code for an access token.
func exchangeCode(ctx context.Context, tokenEndpoint, clientID, clientSecret, code, redirectURL, verifier string) (*TokenData, error) {
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("client_id", clientID)
	v.Set("code", code)
	v.Set("redirect_uri", redirectURL)
	v.Set("code_verifier", verifier)
	if clientSecret != "" {
		v.Set("client_secret", clientSecret)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token endpoint returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("no access_token in response")
	}

	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return &TokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       expiry,
	}, nil
}

func generateVerifier() string {
	b := make([]byte, 48)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:64]
}

func generateS256Challenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
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
