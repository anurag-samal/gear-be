package auth

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	config   *oauth2.Config
}

type GoogleUser struct {
	Subject         string `json:"sub"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Picture         string `json:"picture"`
	EmailVerified   bool   `json:"email_verified"`
}

func NewGoogleProvider(ctx context.Context, clientID string, clientSecret string, redirectURL string) (*GoogleProvider, error) {

	provider, err := oidc.NewProvider(
		ctx,
		"https://accounts.google.com",
	)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			"profile",
			"email",
		},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return &GoogleProvider{
		provider: provider,
		verifier: verifier,
		config:   config,
	}, nil
}

func (g *GoogleProvider) AuthURL(state string, challenge string) string {

	return g.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

func (g *GoogleProvider) ExchangeCode(ctx context.Context,code string,verifier string) (*oauth2.Token, error) {

	return g.config.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)
}

func (g *GoogleProvider) VerifyIDToken(ctx context.Context,token *oauth2.Token) (*GoogleUser, error) {

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id token not found")
	}

	idToken, err := g.verifier.Verify(
		ctx,
		rawIDToken,
	)
	if err != nil {
		return nil, err
	}

	var user GoogleUser

	if err := idToken.Claims(&user); err != nil {
		return nil, err
	}

	return &user, nil
}