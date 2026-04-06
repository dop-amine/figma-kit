package mcpclient

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// TokenData stores the cached OAuth token.
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

func tokenPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "figma-kit", "token.json")
}

// LoadToken reads the cached OAuth token from disk.
func LoadToken() (*TokenData, error) {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return nil, err
	}
	var t TokenData
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// SaveToken writes the OAuth token to disk.
func SaveToken(t *TokenData) error {
	p := tokenPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

// ClearToken removes the cached token file.
func ClearToken() error {
	return os.Remove(tokenPath())
}

// IsValid returns true if the token exists and is not expired.
func (t *TokenData) IsValid() bool {
	if t == nil || t.AccessToken == "" {
		return false
	}
	if t.Expiry.IsZero() {
		return true
	}
	return time.Now().Before(t.Expiry)
}
