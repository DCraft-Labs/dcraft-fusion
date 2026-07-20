package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevModeAllowsRequestWithoutBearerToken(t *testing.T) {
	middleware := NewMiddleware(Config{Mode: "dev"})
	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bootstrap", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected dev request to pass, got %d", response.Code)
	}
}

func TestOIDCModeRequiresBearerToken(t *testing.T) {
	middleware := NewMiddleware(Config{Mode: "oidc", Issuer: "https://issuer.example", Audience: "fusion", JWKSURL: "https://issuer.example/jwks"})
	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bootstrap", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected oidc request to be rejected, got %d", response.Code)
	}
}

func TestOIDCModeAllowsHealthAndMetricsWithoutBearerToken(t *testing.T) {
	middleware := NewMiddleware(Config{Mode: "oidc"})
	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, path := range []string{"/healthz", "/metrics", "/.well-known/openid-configuration", "/oidc/token", "/api/v1/auth/users", "/api/v1/auth/login", "/api/v1/auth/logout"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("expected %s to pass, got %d", path, response.Code)
		}
	}
}
