package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/peterblog/blog/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
)

type ProviderConfig struct {
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier // nil for non-OIDC providers
	IsOIDC       bool
	RequiresPKCE bool
}

var providers map[string]*ProviderConfig

func InitProviders(ctx context.Context) error {
	providers = make(map[string]*ProviderConfig)
	baseURL := os.Getenv("BACKEND_URL")

	// Google
	googleProvider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("google oidc: %w", err)
	}
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	providers["google"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: config.ReadEnv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/google/callback",
			Endpoint:     googleProvider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
		Verifier: googleProvider.Verifier(&oidc.Config{ClientID: googleClientID}),
		IsOIDC:   true,
	}

	// GitHub (OAuth2 only)
	providers["github"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: config.ReadEnv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/github/callback",
			Endpoint:     github.Endpoint,
			Scopes:       []string{"read:user"},
		},
		IsOIDC: false,
	}

	// Facebook (OAuth2 only)
	providers["facebook"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: config.ReadEnv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/facebook/callback",
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"public_profile"},
		},
		IsOIDC: false,
	}

	// LinkedIn (OAuth2 only)
	providers["linkedin"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
			ClientSecret: config.ReadEnv("LINKEDIN_CLIENT_SECRET"),
			RedirectURL:  baseURL + "/auth/linkedin/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
				TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
			},
			Scopes: []string{"openid", "profile"},
		},
		IsOIDC: false,
	}

	// X / Twitter (OAuth2 + PKCE, public client — no secret sent)
	providers["twitter"] = &ProviderConfig{
		OAuth2Config: &oauth2.Config{
			ClientID:    os.Getenv("TWITTER_CLIENT_KEY"),
			RedirectURL: baseURL + "/auth/twitter/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://twitter.com/i/oauth2/authorize",
				TokenURL:  "https://api.twitter.com/2/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
			Scopes: []string{"tweet.read", "users.read"},
		},
		RequiresPKCE: true,
		IsOIDC:       false,
	}

	return nil
}

func GetProvider(name string) (*ProviderConfig, bool) {
	p, ok := providers[name]
	return p, ok
}


type UserInfo struct {
	ProviderUserID string
	DisplayName    string
}

// FetchUserInfo retrieves user identity from the provider after OAuth2 token exchange.
func FetchUserInfo(ctx context.Context, provider string, cfg *ProviderConfig, token *oauth2.Token) (*UserInfo, error) {
	switch provider {
	case "google":
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			return nil, fmt.Errorf("no id_token in response")
		}
		idToken, err := cfg.Verifier.Verify(ctx, rawIDToken)
		if err != nil {
			return nil, fmt.Errorf("id_token verification: %w", err)
		}
		var claims struct {
			Sub  string `json:"sub"`
			Name string `json:"name"`
		}
		if err := idToken.Claims(&claims); err != nil {
			return nil, err
		}
		return &UserInfo{ProviderUserID: claims.Sub, DisplayName: claims.Name}, nil

	case "github":
		return fetchGitHubUser(cfg.OAuth2Config.Client(ctx, token))

	case "facebook":
		return fetchFacebookUser(cfg.OAuth2Config.Client(ctx, token))

	case "linkedin":
		return fetchLinkedInUser(cfg.OAuth2Config.Client(ctx, token))

	case "twitter":
		return fetchXUser(cfg.OAuth2Config.Client(ctx, token))
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
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	name := data.Name
	if name == "" {
		name = data.Login
	}
	return &UserInfo{
		ProviderUserID: fmt.Sprintf("%d", data.ID),
		DisplayName:    name,
	}, nil
}

func fetchFacebookUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &UserInfo{ProviderUserID: data.ID, DisplayName: data.Name}, nil
}

func fetchLinkedInUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://api.linkedin.com/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Sub  string `json:"sub"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &UserInfo{ProviderUserID: data.Sub, DisplayName: data.Name}, nil
}

func fetchXUser(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://api.twitter.com/2/users/me?user.fields=name")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	name := data.Data.Name
	if name == "" {
		name = data.Data.Username
	}
	return &UserInfo{ProviderUserID: data.Data.ID, DisplayName: name}, nil
}

// ExchangeCode exchanges the authorization code for a token.
func ExchangeCode(ctx context.Context, cfg *ProviderConfig, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return cfg.OAuth2Config.Exchange(ctx, code, opts...)
}

// AuthCodeURL returns the provider's authorization URL.
func AuthCodeURL(cfg *ProviderConfig, state string, opts ...oauth2.AuthCodeOption) string {
	return cfg.OAuth2Config.AuthCodeURL(state, opts...)
}

// StorePKCEVerifier saves the PKCE code verifier in the session.
func StorePKCEVerifier(w http.ResponseWriter, r *http.Request, verifier string) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Values["pkce_verifier"] = verifier
	return session.Save(r, w)
}

// GetPKCEVerifier retrieves the PKCE code verifier from the session.
func GetPKCEVerifier(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}
	v, ok := session.Values["pkce_verifier"].(string)
	return v, ok
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
	base := os.Getenv("FRONTEND_URL")
	u, _ := url.Parse(base)
	u.Path = path
	return u.String()
}
