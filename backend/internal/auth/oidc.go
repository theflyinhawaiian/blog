package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
)

type ProviderConfig struct {
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier // nil for non-OIDC providers
	IsOIDC       bool
}

var providers map[string]*ProviderConfig

func InitProviders(ctx context.Context) error {
	providers = make(map[string]*ProviderConfig)
	baseURL := os.Getenv("APP_BASE_URL")

	// Google
	googleProvider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("google oidc: %w", err)
	}
	providers["google"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/google/callback",
			Endpoint:     googleProvider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		Verifier: googleProvider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_CLIENT_ID")}),
		IsOIDC:   true,
	}

	// GitHub (OAuth2 only)
	providers["github"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/github/callback",
			Endpoint:     github.Endpoint,
			Scopes:       []string{"read:user", "user:email"},
		},
		IsOIDC: false,
	}

	// Facebook (OAuth2 only)
	providers["facebook"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/facebook/callback",
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"public_profile", "email"},
		},
		IsOIDC: false,
	}

	// LinkedIn (OAuth2 only)
	providers["linkedin"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
			ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/linkedin/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
				TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
			},
			Scopes: []string{"openid", "profile", "email"},
		},
		IsOIDC: false,
	}

	// Apple (OIDC with JWT client secret)
	appleProvider, err := oidc.NewProvider(ctx, "https://appleid.apple.com")
	if err != nil {
		return fmt.Errorf("apple oidc: %w", err)
	}
	providers["apple"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:    os.Getenv("APPLE_CLIENT_ID"),
			RedirectURL: baseURL + "/auth/apple/callback",
			Endpoint:    appleProvider.Endpoint(),
			Scopes:      []string{oidc.ScopeOpenID, "name", "email"},
		},
		Verifier: appleProvider.Verifier(&oidc.Config{ClientID: os.Getenv("APPLE_CLIENT_ID")}),
		IsOIDC:   true,
	}

	return nil
}

func GetProvider(name string) (*ProviderConfig, bool) {
	p, ok := providers[name]
	return p, ok
}

// GenerateAppleClientSecret builds the JWT needed as Apple's client_secret.
func GenerateAppleClientSecret() (string, error) {
	keyPEM := os.Getenv("APPLE_CLIENT_SECRET") // the .p8 private key PEM
	if keyPEM == "" {
		return "", fmt.Errorf("APPLE_CLIENT_SECRET not set")
	}

	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse apple private key: %w", err)
	}

	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("apple key is not ECDSA")
	}

	header := base64.RawURLEncoding.EncodeToString(mustMarshalJSON(map[string]string{
		"alg": "ES256",
		"kid": os.Getenv("APPLE_KEY_ID"),
	}))
	now := time.Now()
	payload := base64.RawURLEncoding.EncodeToString(mustMarshalJSON(map[string]interface{}{
		"iss": os.Getenv("APPLE_TEAM_ID"),
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": os.Getenv("APPLE_CLIENT_ID"),
	}))

	sigInput := header + "." + payload
	hash := sha256.Sum256([]byte(sigInput))

	r, s, err := ecdsa.Sign(rand.Reader, ecKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("signing: %w", err)
	}

	// IEEE P1363 format: r || s, each 32 bytes
	sig := make([]byte, 64)
	rBytes := paddedBigInt(r, 32)
	sBytes := paddedBigInt(s, 32)
	copy(sig[:32], rBytes)
	copy(sig[32:], sBytes)

	return sigInput + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func paddedBigInt(n *big.Int, size int) []byte {
	b := n.Bytes()
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

func mustMarshalJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

type UserInfo struct {
	ProviderUserID string
	DisplayName    string
	Email          *string
}

// FetchUserInfo retrieves user identity from the provider after OAuth2 token exchange.
func FetchUserInfo(ctx context.Context, provider string, cfg *ProviderConfig, token *oauth2.Token) (*UserInfo, error) {
	switch provider {
	case "google", "apple":
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			return nil, fmt.Errorf("no id_token in response")
		}
		idToken, err := cfg.Verifier.Verify(ctx, rawIDToken)
		if err != nil {
			return nil, fmt.Errorf("id_token verification: %w", err)
		}
		var claims struct {
			Sub   string `json:"sub"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			return nil, err
		}
		info := &UserInfo{ProviderUserID: claims.Sub, DisplayName: claims.Name}
		if claims.Email != "" {
			info.Email = &claims.Email
		}
		return info, nil

	case "github":
		return fetchGitHubUser(cfg.OAuth2Config.Client(ctx, token))

	case "facebook":
		return fetchFacebookUser(cfg.OAuth2Config.Client(ctx, token))

	case "linkedin":
		return fetchLinkedInUser(cfg.OAuth2Config.Client(ctx, token))
	}

	return nil, fmt.Errorf("unknown provider: %s", provider)
}

func fetchGitHubUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	name := data.Name
	if name == "" {
		name = data.Login
	}
	info := &UserInfo{
		ProviderUserID: fmt.Sprintf("%d", data.ID),
		DisplayName:    name,
	}
	if data.Email != "" {
		info.Email = &data.Email
	}
	return info, nil
}

func fetchFacebookUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	info := &UserInfo{ProviderUserID: data.ID, DisplayName: data.Name}
	if data.Email != "" {
		info.Email = &data.Email
	}
	return info, nil
}

func fetchLinkedInUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://api.linkedin.com/v2/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ID             string `json:"id"`
		LocalizedFName string `json:"localizedFirstName"`
		LocalizedLName string `json:"localizedLastName"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(data.LocalizedFName + " " + data.LocalizedLName)
	return &UserInfo{ProviderUserID: data.ID, DisplayName: name}, nil
}

// ExchangeCode exchanges the authorization code for a token.
func ExchangeCode(ctx context.Context, cfg *ProviderConfig, code string, extraSecret ...string) (*oauth2.Token, error) {
	oauthCfg := *cfg.OAuth2Config
	if len(extraSecret) > 0 {
		oauthCfg.ClientSecret = extraSecret[0]
	}
	return oauthCfg.Exchange(ctx, code)
}

// AuthCodeURL returns the provider's authorization URL.
func AuthCodeURL(cfg *ProviderConfig, state string) string {
	return cfg.OAuth2Config.AuthCodeURL(state)
}

// OAuthStateParam generates a CSRF state token.
func OAuthStateParam() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// StoreOAuthState saves state in session for CSRF check.
func StoreOAuthState(w http.ResponseWriter, r *http.Request, state string) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Values["oauth_state"] = state
	return session.Save(r, w)
}

// ValidateOAuthState checks state matches.
func ValidateOAuthState(r *http.Request, state string) bool {
	session, err := GetSession(r)
	if err != nil {
		return false
	}
	stored, ok := session.Values["oauth_state"].(string)
	return ok && stored == state
}

// RedirectToFrontend builds a frontend redirect URL.
func RedirectToFrontend(path string) string {
	base := os.Getenv("APP_BASE_URL")
	u, _ := url.Parse(base)
	u.Path = path
	return u.String()
}
