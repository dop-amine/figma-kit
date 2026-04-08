package restapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	srv := httptest.NewServer(handler)
	c := &Client{
		token:      "test-pat",
		httpClient: srv.Client(),
	}
	return c, srv
}

func TestGetTeamComponents(t *testing.T) {
	want := ComponentsResponse{
		Meta: ComponentsMeta{
			Components: []Component{
				{Key: "abc123", Name: "Button", Description: "Primary button"},
			},
		},
	}
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Figma-Token") != "test-pat" {
			t.Error("missing auth header")
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/teams/999/components") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(want)
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	resp, err := c.GetTeamComponents("999", 10, "")
	if err != nil {
		t.Fatalf("GetTeamComponents error: %v", err)
	}
	if len(resp.Meta.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(resp.Meta.Components))
	}
	if resp.Meta.Components[0].Key != "abc123" {
		t.Errorf("key = %q, want abc123", resp.Meta.Components[0].Key)
	}
}

func TestGetFileComponents(t *testing.T) {
	want := ComponentsResponse{
		Meta: ComponentsMeta{
			Components: []Component{
				{Key: "file-comp", Name: "Card"},
			},
		},
	}
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/files/fileXYZ/components") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(want)
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	resp, err := c.GetFileComponents("fileXYZ")
	if err != nil {
		t.Fatalf("GetFileComponents error: %v", err)
	}
	if resp.Meta.Components[0].Name != "Card" {
		t.Errorf("name = %q, want Card", resp.Meta.Components[0].Name)
	}
}

func TestGetComponent(t *testing.T) {
	want := SingleComponentResponse{
		Meta: Component{Key: "k1", Name: "Icon", Description: "An icon"},
	}
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/components/k1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(want)
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	comp, err := c.GetComponent("k1")
	if err != nil {
		t.Fatalf("GetComponent error: %v", err)
	}
	if comp.Name != "Icon" {
		t.Errorf("name = %q, want Icon", comp.Name)
	}
}

func TestGetStyle(t *testing.T) {
	want := SingleStyleResponse{
		Meta: Style{Key: "s1", Name: "Primary Fill", StyleType: "FILL"},
	}
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/styles/s1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(want)
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	style, err := c.GetStyle("s1")
	if err != nil {
		t.Fatalf("GetStyle error: %v", err)
	}
	if style.StyleType != "FILL" {
		t.Errorf("styleType = %q, want FILL", style.StyleType)
	}
}

func TestGetTeamStyles(t *testing.T) {
	want := StylesResponse{
		Meta: StylesMeta{
			Styles: []Style{
				{Key: "ts1", Name: "Body Text", StyleType: "TEXT"},
			},
		},
	}
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/teams/42/styles") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(want)
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	resp, err := c.GetTeamStyles("42", 0, "")
	if err != nil {
		t.Fatalf("GetTeamStyles error: %v", err)
	}
	if len(resp.Meta.Styles) != 1 || resp.Meta.Styles[0].Key != "ts1" {
		t.Errorf("unexpected styles result: %+v", resp.Meta.Styles)
	}
}

func TestPagination(t *testing.T) {
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page_size") != "5" {
			t.Errorf("page_size = %q, want 5", q.Get("page_size"))
		}
		if q.Get("after") != "cursor123" {
			t.Errorf("after = %q, want cursor123", q.Get("after"))
		}
		json.NewEncoder(w).Encode(ComponentsResponse{})
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	_, err := c.GetTeamComponents("1", 5, "cursor123")
	if err != nil {
		t.Fatalf("pagination error: %v", err)
	}
}

func TestHTTPError(t *testing.T) {
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"status":403,"err":"Forbidden"}`))
	})
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	_, err := c.GetComponent("bad")
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should contain 403: %v", err)
	}
}
