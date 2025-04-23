package oauth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rhajizada/gazette/internal/config"
	"golang.org/x/oauth2"
)

func GetConfig(cfg *config.OAuthConfig) (*oauth2.Config, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		RedirectURL:  cfg.RedirectURL,
	}
	return oauthCfg, nil
}

func GetVerifier(cfg *config.OAuthConfig) (*oidc.IDTokenVerifier, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})
	return verifier, nil
}
