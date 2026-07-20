package auth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProviderIssuesTokenAndMiddlewareAcceptsIt(t *testing.T) {
	provider, err := NewProvider(Config{
		Mode:     "oidc",
		Issuer:   "http://127.0.0.1:30173",
		Audience: "dcraft-fusion",
		JWKSURL:  "",
	})
	if err != nil {
		t.Fatalf("provider init failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", provider.Login)
	mux.HandleFunc("GET /oidc/authorize", provider.Authorize)
	mux.HandleFunc("POST /oidc/token", provider.Token)
	mux.HandleFunc("GET /oidc/jwks", provider.JWKS)
	server := httptest.NewServer(mux)
	defer server.Close()
	provider.config.Issuer = server.URL

	middleware := NewMiddleware(Config{
		Mode:     "oidc",
		Issuer:   server.URL,
		Audience: "dcraft-fusion",
		JWKSURL:  server.URL + "/oidc/jwks",
	})
	protected := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Actor", r.Header.Get("X-Fusion-Actor-Id"))
		w.Header().Set("X-Organization", r.Header.Get("X-Fusion-Organization-Id"))
		w.WriteHeader(http.StatusNoContent)
	}))

	loginBody := []byte(`{"email":"founder@dcraftlabs.com","password":"changeme-founder"}`)
	loginResponse := httptest.NewRecorder()
	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginRequest.Header.Set("Content-Type", "application/json")
	provider.Login(loginResponse, loginRequest)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("expected login success, got %d", loginResponse.Code)
	}
	loginCookie := loginResponse.Result().Cookies()
	if len(loginCookie) == 0 {
		t.Fatal("expected login session cookie")
	}

	authorizeResponse := httptest.NewRecorder()
	authorizeRequest := httptest.NewRequest(http.MethodGet, "/oidc/authorize?client_id=fusion-web&redirect_uri=http://localhost/auth/callback&response_type=code&state=test-state", nil)
	authorizeRequest.AddCookie(loginCookie[0])
	provider.Authorize(authorizeResponse, authorizeRequest)
	if authorizeResponse.Code != http.StatusFound {
		t.Fatalf("expected authorize redirect, got %d", authorizeResponse.Code)
	}
	location, err := url.Parse(authorizeResponse.Header().Get("Location"))
	if err != nil {
		t.Fatalf("authorize location invalid: %v", err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatal("authorization code missing")
	}

	form := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {"fusion-web"},
		"redirect_uri": {"http://localhost/auth/callback"},
		"code":         {code},
	}
	tokenResponse := httptest.NewRecorder()
	tokenRequest := httptest.NewRequest(http.MethodPost, "/oidc/token", strings.NewReader(form.Encode()))
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	provider.Token(tokenResponse, tokenRequest)
	if tokenResponse.Code != http.StatusOK {
		t.Fatalf("expected token response 200, got %d with body %s", tokenResponse.Code, tokenResponse.Body.String())
	}
	body, err := io.ReadAll(tokenResponse.Body)
	if err != nil {
		t.Fatalf("failed to read token response: %v", err)
	}
	if !strings.Contains(string(body), "access_token") {
		t.Fatalf("token response missing access token: %s", string(body))
	}
	token := extractJSONField(string(body), "access_token")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/bootstrap", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Fusion-Organization-Id", "forged-org")
	response := httptest.NewRecorder()

	protected.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected protected request to pass, got %d", response.Code)
	}
	if response.Header().Get("X-Actor") != "user-founder" {
		t.Fatalf("expected actor from token, got %q", response.Header().Get("X-Actor"))
	}
	if response.Header().Get("X-Organization") != "org-b2b2c-1" {
		t.Fatalf("expected organization from token, got %q", response.Header().Get("X-Organization"))
	}
}

func extractJSONField(body string, field string) string {
	prefix := `"` + field + `":"`
	start := strings.Index(body, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(body[start:], `"`)
	if end == -1 {
		return ""
	}
	return body[start : start+end]
}
