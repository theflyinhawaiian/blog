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
	"strings"

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
	baseURL := os.Getenv("BACKEND_URL")

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

	return nil
}

func GetProvider(name string) (*ProviderConfig, bool) {
	p, ok := providers[name]
	return p, ok
}


type UserInfo struct {
	ProviderUserID string
	DisplayName    string
	Email          *string
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
	base := os.Getenv("FRONTEND_URL")
	u, _ := url.Parse(base)
	u.Path = path
	return u.String()
}
