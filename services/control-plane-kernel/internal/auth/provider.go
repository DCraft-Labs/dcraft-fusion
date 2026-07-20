package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const webClientID = "fusion-web"
const sessionCookieName = "fusion_oidc_session"

type SeededUser struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Label          string `json:"label"`
	Scope          string `json:"scope"`
	Password       string `json:"-"`
	PasswordHint   string `json:"passwordHint,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
	TenantID       string `json:"tenantId,omitempty"`
	ProjectID      string `json:"projectId,omitempty"`
}

type authorizationCode struct {
	ClientID    string
	RedirectURI string
	User        SeededUser
	ExpiresAt   time.Time
}

type loginSession struct {
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Provider struct {
	config Config
	keyID  string
	key    *rsa.PrivateKey

	mu       sync.Mutex
	codes    map[string]authorizationCode
	sessions map[string]loginSession
	users    []SeededUser
}

func NewProvider(config Config) (*Provider, error) {
	var key *rsa.PrivateKey
	var err error
	if strings.TrimSpace(config.PrivateKeyPEM) != "" {
		key, err = parseRSAPrivateKey(config.PrivateKeyPEM)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
	}
	sum := sha256.Sum256(key.N.Bytes())
	return &Provider{
		config: config,
		keyID:  hex.EncodeToString(sum[:8]),
		key:    key,
		codes:  map[string]authorizationCode{},
		users: seededUsersFromEnv(),
		sessions: map[string]loginSession{},
	}, nil
}

func (provider *Provider) Users(w http.ResponseWriter, _ *http.Request) {
	publicUsers := make([]map[string]any, 0, len(provider.users))
	for _, user := range provider.users {
		publicUsers = append(publicUsers, map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"label": user.Label,
			"scope": user.Scope,
		})
	}
	writeOIDCJSON(w, http.StatusOK, map[string]any{"users": publicUsers})
}

func (provider *Provider) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid login request", http.StatusBadRequest)
		return
	}
	user, ok := provider.userByEmail(request.Email)
	if !ok || user.Password != request.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	sessionID := provider.issueSession(user)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((8 * time.Hour).Seconds()),
	})

	// Community / password mode: return a signed JWT directly (no OIDC redirect).
	if strings.EqualFold(provider.config.Mode, "password") || strings.EqualFold(provider.config.Mode, "local") {
		token, expiresAt, err := provider.signToken(user)
		if err != nil {
			http.Error(w, "failed to sign token", http.StatusInternalServerError)
			return
		}
		writeOIDCJSON(w, http.StatusOK, map[string]any{
			"ok":           true,
			"email":        user.Email,
			"scope":        user.Scope,
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   int(time.Until(expiresAt).Seconds()),
		})
		return
	}

	writeOIDCJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"email": user.Email,
		"scope": user.Scope,
	})
}

func seededUsersFromEnv() []SeededUser {
	founderPassword := firstNonEmpty(os.Getenv("FUSION_SEED_FOUNDER_PASSWORD"), "changeme-founder")
	superadminPassword := firstNonEmpty(os.Getenv("FUSION_SEED_SUPERADMIN_PASSWORD"), "changeme-superadmin")
	return []SeededUser{
		{
			ID:             "user-founder",
			Email:          firstNonEmpty(os.Getenv("FUSION_SEED_FOUNDER_EMAIL"), "founder@dcraftlabs.com"),
			Label:          "Founder Admin",
			Scope:          "tenant",
			Password:       founderPassword,
			OrganizationID: "org-b2b2c-1",
			TenantID:       "tenant-brand-a",
			ProjectID:      "project-prod",
		},
		{
			ID:       "user-superadmin",
			Email:    firstNonEmpty(os.Getenv("FUSION_SEED_SUPERADMIN_EMAIL"), "superadmin@dcraftlabs.com"),
			Label:    "Platform Superadmin",
			Scope:    "platform",
			Password: superadminPassword,
		},
	}
}

func (provider *Provider) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && strings.TrimSpace(cookie.Value) != "" {
		provider.deleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	writeOIDCJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (provider *Provider) OpenIDConfiguration(w http.ResponseWriter, _ *http.Request) {
	writeOIDCJSON(w, http.StatusOK, map[string]any{
		"issuer":                                provider.config.Issuer,
		"authorization_endpoint":                provider.config.Issuer + "/oidc/authorize",
		"token_endpoint":                        provider.config.Issuer + "/oidc/token",
		"jwks_uri":                              provider.config.Issuer + "/oidc/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"token_endpoint_auth_methods_supported": []string{"none"},
	})
}

func (provider *Provider) JWKS(w http.ResponseWriter, _ *http.Request) {
	exponent := big.NewInt(int64(provider.key.PublicKey.E)).Bytes()
	writeOIDCJSON(w, http.StatusOK, map[string]any{
		"keys": []map[string]string{
			{
				"kid": provider.keyID,
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"n":   base64.RawURLEncoding.EncodeToString(provider.key.PublicKey.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(exponent),
			},
		},
	})
}

func (provider *Provider) Authorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("response_type") != "code" {
		http.Error(w, "unsupported response_type", http.StatusBadRequest)
		return
	}
	clientID := strings.TrimSpace(query.Get("client_id"))
	redirectURI := strings.TrimSpace(query.Get("redirect_uri"))
	if clientID != webClientID || redirectURI == "" {
		http.Error(w, "invalid authorize request", http.StatusBadRequest)
		return
	}
	user, ok := provider.authorizedUser(r)
	if !ok {
		http.Error(w, "login required", http.StatusUnauthorized)
		return
	}
	if fallbackUserID := firstNonEmpty(query.Get("user_id"), query.Get("login_hint")); fallbackUserID != "" {
		user, ok = provider.userByID(fallbackUserID)
	}
	if !ok {
		http.Error(w, "unknown user", http.StatusBadRequest)
		return
	}
	code := provider.issueCode(authorizationCode{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		User:        user,
		ExpiresAt:   time.Now().UTC().Add(5 * time.Minute),
	})
	location, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}
	params := location.Query()
	params.Set("code", code)
	if state := strings.TrimSpace(query.Get("state")); state != "" {
		params.Set("state", state)
	}
	location.RawQuery = params.Encode()
	http.Redirect(w, r, location.String(), http.StatusFound)
}

func (provider *Provider) Token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid token request", http.StatusBadRequest)
		return
	}
	if r.Form.Get("grant_type") != "authorization_code" {
		http.Error(w, "unsupported grant_type", http.StatusBadRequest)
		return
	}
	code := strings.TrimSpace(r.Form.Get("code"))
	clientID := strings.TrimSpace(r.Form.Get("client_id"))
	redirectURI := strings.TrimSpace(r.Form.Get("redirect_uri"))
	entry, err := provider.consumeCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if entry.ClientID != clientID || entry.ClientID != webClientID || entry.RedirectURI != redirectURI {
		http.Error(w, "authorization code mismatch", http.StatusBadRequest)
		return
	}
	token, expiresAt, err := provider.signToken(entry.User)
	if err != nil {
		http.Error(w, "failed to sign token", http.StatusInternalServerError)
		return
	}
	writeOIDCJSON(w, http.StatusOK, map[string]any{
		"access_token": token,
		"id_token":     token,
		"token_type":   "Bearer",
		"expires_in":   int(time.Until(expiresAt).Seconds()),
		"scope":        "openid profile email",
	})
}

func (provider *Provider) signToken(user SeededUser) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(60 * time.Minute)
	header := map[string]string{
		"alg": "RS256",
		"kid": provider.keyID,
		"typ": "JWT",
	}
	claims := map[string]any{
		"iss":             provider.config.Issuer,
		"sub":             user.ID,
		"aud":             provider.config.Audience,
		"exp":             expiresAt.Unix(),
		"iat":             time.Now().UTC().Unix(),
		"actor_id":        user.ID,
		"email":           user.Email,
		"name":            user.Label,
		"scope":           user.Scope,
		"organization_id": user.OrganizationID,
		"tenant_id":       user.TenantID,
		"project_id":      user.ProjectID,
	}
	encodedHeader, err := encodeJWTPart(header)
	if err != nil {
		return "", time.Time{}, err
	}
	encodedClaims, err := encodeJWTPart(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	signed := encodedHeader + "." + encodedClaims
	digest := sha256.Sum256([]byte(signed))
	signature, err := rsa.SignPKCS1v15(rand.Reader, provider.key, crypto.SHA256, digest[:])
	if err != nil {
		return "", time.Time{}, err
	}
	return signed + "." + base64.RawURLEncoding.EncodeToString(signature), expiresAt, nil
}

func encodeJWTPart(value any) (string, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(body), nil
}

func (provider *Provider) issueCode(code authorizationCode) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", code.User.ID, code.RedirectURI, time.Now().UnixNano())))
	issued := base64.RawURLEncoding.EncodeToString(sum[:])
	if err := provider.persistCode(issued, code); err == nil {
		return issued
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.codes[issued] = code
	return issued
}

func (provider *Provider) consumeCode(code string) (authorizationCode, error) {
	if entry, err := provider.loadCode(code); err == nil {
		return entry, nil
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	entry, ok := provider.codes[code]
	if !ok {
		return authorizationCode{}, errors.New("authorization code not found")
	}
	delete(provider.codes, code)
	if time.Now().UTC().After(entry.ExpiresAt) {
		return authorizationCode{}, errors.New("authorization code expired")
	}
	return entry, nil
}

func (provider *Provider) userByID(id string) (SeededUser, bool) {
	for _, user := range provider.users {
		if user.ID == strings.TrimSpace(id) {
			return user, true
		}
	}
	return SeededUser{}, false
}

func (provider *Provider) userByEmail(email string) (SeededUser, bool) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	for _, user := range provider.users {
		if strings.ToLower(user.Email) == normalized {
			return user, true
		}
	}
	return SeededUser{}, false
}

func writeOIDCJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (provider *Provider) persistCode(code string, entry authorizationCode) error {
	if strings.TrimSpace(provider.config.RedisAddr) == "" {
		return errors.New("redis unavailable")
	}
	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	ttl := int(time.Until(entry.ExpiresAt).Seconds())
	if ttl <= 0 {
		ttl = 1
	}
	return redisSetJSON(provider.config.RedisAddr, redisCodeKey(code), string(body), ttl)
}

func (provider *Provider) loadCode(code string) (authorizationCode, error) {
	if strings.TrimSpace(provider.config.RedisAddr) == "" {
		return authorizationCode{}, errors.New("redis unavailable")
	}
	body, err := redisGet(provider.config.RedisAddr, redisCodeKey(code))
	if err != nil {
		return authorizationCode{}, err
	}
	_ = redisDelete(provider.config.RedisAddr, redisCodeKey(code))
	var entry authorizationCode
	if err := json.Unmarshal([]byte(body), &entry); err != nil {
		return authorizationCode{}, err
	}
	if time.Now().UTC().After(entry.ExpiresAt) {
		return authorizationCode{}, errors.New("authorization code expired")
	}
	return entry, nil
}

func (provider *Provider) issueSession(user SeededUser) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("session:%s:%d", user.ID, time.Now().UnixNano())))
	sessionID := base64.RawURLEncoding.EncodeToString(sum[:])
	entry := loginSession{UserID: user.ID, ExpiresAt: time.Now().UTC().Add(8 * time.Hour)}
	if err := provider.persistSession(sessionID, entry); err == nil {
		return sessionID
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.sessions[sessionID] = entry
	return sessionID
}

func (provider *Provider) authorizedUser(r *http.Request) (SeededUser, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return SeededUser{}, false
	}
	entry, err := provider.loadSession(cookie.Value)
	if err != nil {
		return SeededUser{}, false
	}
	return provider.userByID(entry.UserID)
}

func (provider *Provider) persistSession(sessionID string, entry loginSession) error {
	if strings.TrimSpace(provider.config.RedisAddr) == "" {
		return errors.New("redis unavailable")
	}
	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	ttl := int(time.Until(entry.ExpiresAt).Seconds())
	if ttl <= 0 {
		ttl = 1
	}
	return redisSetJSON(provider.config.RedisAddr, redisSessionKey(sessionID), string(body), ttl)
}

func (provider *Provider) loadSession(sessionID string) (loginSession, error) {
	if strings.TrimSpace(provider.config.RedisAddr) != "" {
		body, err := redisGet(provider.config.RedisAddr, redisSessionKey(sessionID))
		if err == nil {
			var entry loginSession
			if err := json.Unmarshal([]byte(body), &entry); err != nil {
				return loginSession{}, err
			}
			if time.Now().UTC().After(entry.ExpiresAt) {
				return loginSession{}, errors.New("login session expired")
			}
			return entry, nil
		}
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	entry, ok := provider.sessions[sessionID]
	if !ok {
		return loginSession{}, errors.New("login session not found")
	}
	if time.Now().UTC().After(entry.ExpiresAt) {
		delete(provider.sessions, sessionID)
		return loginSession{}, errors.New("login session expired")
	}
	return entry, nil
}

func (provider *Provider) deleteSession(sessionID string) {
	if strings.TrimSpace(provider.config.RedisAddr) != "" {
		_ = redisDelete(provider.config.RedisAddr, redisSessionKey(sessionID))
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	delete(provider.sessions, sessionID)
}
