package auth

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
)

type Config struct {
	Mode          string
	Issuer        string
	Audience      string
	JWKSURL       string
	RedisAddr     string
	PrivateKeyPEM string
}

func FromEnv() Config {
	return Config{
		Mode:          firstNonEmpty(os.Getenv("FUSION_AUTH_MODE"), "dev"),
		Issuer:        os.Getenv("FUSION_OIDC_ISSUER"),
		Audience:      os.Getenv("FUSION_OIDC_AUDIENCE"),
		JWKSURL:       os.Getenv("FUSION_OIDC_JWKS_URL"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		PrivateKeyPEM: os.Getenv("FUSION_OIDC_PRIVATE_KEY_PEM"),
	}
}

type Middleware struct {
	config Config
	keys   *jwksCache
}

func NewMiddleware(config Config) *Middleware {
	return &Middleware{config: config, keys: &jwksCache{url: config.JWKSURL}}
}

func (middleware *Middleware) Wrap(next http.Handler) http.Handler {
	mode := strings.ToLower(strings.TrimSpace(middleware.config.Mode))
	if mode == "dev" || mode == "" {
		return next
	}
	// password / local / oidc all require bearer tokens on protected routes.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		claims, err := middleware.verifyBearer(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		r.Header.Del(fusioncontext.ActorIDHeader)
		r.Header.Del(fusioncontext.OrganizationIDHeader)
		r.Header.Del(fusioncontext.TenantIDHeader)
		r.Header.Del(fusioncontext.ProjectIDHeader)
		r.Header.Set(fusioncontext.ActorIDHeader, firstNonEmpty(claims.Subject, claims.ActorID))
		if claims.OrganizationID != "" {
			r.Header.Set(fusioncontext.OrganizationIDHeader, claims.OrganizationID)
		}
		if claims.TenantID != "" {
			r.Header.Set(fusioncontext.TenantIDHeader, claims.TenantID)
		}
		if claims.ProjectID != "" {
			r.Header.Set(fusioncontext.ProjectIDHeader, claims.ProjectID)
		}
		if r.Header.Get(fusioncontext.CorrelationIDHeader) == "" {
			r.Header.Set(fusioncontext.CorrelationIDHeader, "oidc-"+time.Now().UTC().Format("20060102150405.000000000"))
		}
		next.ServeHTTP(w, r)
	})
}

func isPublicPath(path string) bool {
	return path == "/healthz" ||
		path == "/metrics" ||
		path == "/.well-known/openid-configuration" ||
		strings.HasPrefix(path, "/oidc/") ||
		path == "/api/v1/auth/users" ||
		path == "/api/v1/auth/login" ||
		path == "/api/v1/auth/logout"
}

type tokenClaims struct {
	Issuer         string          `json:"iss"`
	Subject        string          `json:"sub"`
	Audience       json.RawMessage `json:"aud"`
	ExpiresAt      int64           `json:"exp"`
	ActorID        string          `json:"actor_id"`
	OrganizationID string          `json:"organization_id"`
	TenantID       string          `json:"tenant_id"`
	ProjectID      string          `json:"project_id"`
}

func (middleware *Middleware) verifyBearer(ctx context.Context, authorization string) (tokenClaims, error) {
	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == authorization || token == "" {
		return tokenClaims{}, errors.New("missing bearer token")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenClaims{}, errors.New("invalid bearer token")
	}
	var header struct {
		Algorithm string `json:"alg"`
		KeyID     string `json:"kid"`
	}
	if err := decodeJSON(parts[0], &header); err != nil {
		return tokenClaims{}, errors.New("invalid token header")
	}
	if header.Algorithm != "RS256" {
		return tokenClaims{}, errors.New("unsupported token algorithm")
	}
	var claims tokenClaims
	if err := decodeJSON(parts[1], &claims); err != nil {
		return tokenClaims{}, errors.New("invalid token claims")
	}
	if claims.Issuer != middleware.config.Issuer {
		return tokenClaims{}, errors.New("token issuer mismatch")
	}
	if !claims.hasAudience(middleware.config.Audience) {
		return tokenClaims{}, errors.New("token audience mismatch")
	}
	if claims.ExpiresAt <= time.Now().Unix() {
		return tokenClaims{}, errors.New("token expired")
	}
	key, err := middleware.keys.key(ctx, header.KeyID)
	if err != nil {
		return tokenClaims{}, err
	}
	signed := parts[0] + "." + parts[1]
	digest := sha256.Sum256([]byte(signed))
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return tokenClaims{}, errors.New("invalid token signature")
	}
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature); err != nil {
		return tokenClaims{}, errors.New("token signature verification failed")
	}
	return claims, nil
}

func (claims tokenClaims) hasAudience(expected string) bool {
	var single string
	if err := json.Unmarshal(claims.Audience, &single); err == nil {
		return single == expected
	}
	var many []string
	if err := json.Unmarshal(claims.Audience, &many); err != nil {
		return false
	}
	for _, audience := range many {
		if audience == expected {
			return true
		}
	}
	return false
}

type jwksCache struct {
	mu        sync.Mutex
	url       string
	keys      map[string]*rsa.PublicKey
	expiresAt time.Time
}

func (cache *jwksCache) key(ctx context.Context, keyID string) (*rsa.PublicKey, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.keys != nil && time.Now().Before(cache.expiresAt) {
		if key := cache.keys[keyID]; key != nil {
			return key, nil
		}
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cache.url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var body struct {
		Keys []struct {
			KeyID string `json:"kid"`
			KTY   string `json:"kty"`
			N     string `json:"n"`
			E     string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, err
	}
	cache.keys = map[string]*rsa.PublicKey{}
	for _, jwk := range body.Keys {
		if jwk.KTY != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}
		cache.keys[jwk.KeyID] = &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: int(new(big.Int).SetBytes(eBytes).Int64())}
	}
	cache.expiresAt = time.Now().Add(5 * time.Minute)
	if key := cache.keys[keyID]; key != nil {
		return key, nil
	}
	return nil, errors.New("jwks key not found")
}

func decodeJSON(segment string, value any) error {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func parseRSAPrivateKey(pemValue string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemValue))
	if block == nil {
		return nil, errors.New("invalid oidc private key pem")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("oidc private key is not rsa")
		}
		return rsaKey, nil
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func redisSetJSON(addr string, key string, value string, ttlSeconds int) error {
	if strings.TrimSpace(addr) == "" {
		return errors.New("redis unavailable")
	}
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	command := []string{"SET", key, value, "EX", strconv.Itoa(ttlSeconds)}
	if _, err := conn.Write([]byte(resp(command))); err != nil {
		return err
	}
	buffer := make([]byte, 16)
	_, err = conn.Read(buffer)
	return err
}

func redisGet(addr string, key string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(resp([]string{"GET", key}))); err != nil {
		return "", err
	}
	buffer := make([]byte, 65536)
	count, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	if count == 0 || buffer[0] != '$' {
		return "", errors.New("redis get failed")
	}
	payload := string(buffer[:count])
	if strings.HasPrefix(payload, "$-1") {
		return "", errors.New("redis key not found")
	}
	parts := strings.SplitN(payload, "\r\n", 3)
	if len(parts) < 3 {
		return "", errors.New("redis get malformed response")
	}
	return parts[1], nil
}

func redisDelete(addr string, key string) error {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(resp([]string{"DEL", key}))); err != nil {
		return err
	}
	buffer := make([]byte, 16)
	_, err = conn.Read(buffer)
	return err
}

func resp(parts []string) string {
	var builder strings.Builder
	builder.WriteString("*")
	builder.WriteString(strconv.Itoa(len(parts)))
	builder.WriteString("\r\n")
	for _, part := range parts {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(len(part)))
		builder.WriteString("\r\n")
		builder.WriteString(part)
		builder.WriteString("\r\n")
	}
	return builder.String()
}

func redisCodeKey(code string) string {
	return fmt.Sprintf("fusion:oidc:code:%s", code)
}

func redisSessionKey(sessionID string) string {
	return fmt.Sprintf("fusion:oidc:session:%s", sessionID)
}
