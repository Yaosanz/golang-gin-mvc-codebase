package provider

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"go-starter-app/config"
	"log"
	"os"
)

type OidcProvider struct {
	verifier *oidc.IDTokenVerifier
	issuer   string
	logger   *log.Logger
}

type IOidcProvider interface {
	VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error)
	ParseClaims(idToken *oidc.IDToken, claims interface{}) error
	VerifyAndParseClaims(ctx context.Context, token string, claims interface{}) error
	Logger() *log.Logger
}

// NewOidcProvider initialize oidc provider
func NewOidcProvider(context context.Context, config *config.Config) (*OidcProvider, error) {
	provider, err := oidc.NewProvider(context, config.Oidc().Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID:          config.Oidc().ClientID,
		SkipClientIDCheck: config.Oidc().SkipClientIDCheck, // Skipping default client ID check
	}

	return &OidcProvider{
		verifier: provider.Verifier(oidcConfig),
		issuer:   config.Oidc().Issuer,
		logger:   log.New(os.Stdout, "[OIDC-VERIFIER] ", log.LstdFlags|log.Lshortfile),
	}, nil
}

// VerifyToken verify JWT token (access_token) and returns IDToken
func (v *OidcProvider) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		v.logger.Printf("failed to verify token: %v", err)
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return idToken, nil
}

// ParseClaims parse claim from verified token
func (v *OidcProvider) ParseClaims(idToken *oidc.IDToken, claims interface{}) error {
	// decode id
	if err := idToken.Claims(claims); err != nil {
		v.logger.Printf("failed to parse claims: %v", err)
		return fmt.Errorf("failed to parse claims: %w", err)
	}

	return nil
}

// VerifyAndParseClaims verify and parse claims
func (v *OidcProvider) VerifyAndParseClaims(ctx context.Context, token string, claims interface{}) error {
	idToken, err := v.VerifyToken(ctx, token)
	if err != nil {
		return err
	}

	if err := idToken.Claims(claims); err != nil {
		v.logger.Printf("failed to parse claims: %v", err)
		return fmt.Errorf("failed to parse claims: %w", err)
	}

	return nil
}

func (v *OidcProvider) Logger() *log.Logger {
	return v.logger
}
