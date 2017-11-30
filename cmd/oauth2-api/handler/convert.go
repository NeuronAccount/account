package handler

import (
	api "github.com/NeuronGroup/Account/api/oauth2/gen/models"
	"github.com/NeuronGroup/Account/api/oauth2/gen/restapi/operations"
	"github.com/NeuronGroup/Account/models"
)

func toTokenRequest(p *operations.OAuth2TokenParams) (r *models.OAuth2TokenRequest) {
	if p == nil {
		return nil
	}

	r = &models.OAuth2TokenRequest{}
	r.GrantType = p.GrantType
	if p.Code != nil {
		r.AuthorizationCode = *p.Code
	}
	if p.ResponseType != nil {
		r.ResponseType = *p.ResponseType
	}
	if p.RedirectURI != nil {
		r.RedirectURI = *p.RedirectURI
	}
	if p.ClientID != nil {
		r.ClientID = *p.ClientID
	}
	if p.Scope != nil {
		r.Scope = *p.Scope
	}
	if p.RefreshToken != nil {
		r.RefreshToken = *p.RefreshToken
	}

	return r
}

func fromTokenResponse(p *models.OAuth2AccessToken) (r *api.AccessToken) {
	if p == nil {
		return nil
	}

	r = &api.AccessToken{}
	r.TokenType = p.TokenType
	r.AccessToken = p.AccessToken
	r.ExpiresIn = p.ExpiresIn
	r.RefreshToken = p.RefreshToken
	r.Scope = p.Scope

	return r
}
