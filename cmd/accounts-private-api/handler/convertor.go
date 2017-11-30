package handler

import (
	api "github.com/NeuronGroup/Account/api/private/gen/models"
	"github.com/NeuronGroup/Account/models"
)

func toOAuth2Param(p *api.OAuth2AuthorizeParams) (r *models.OAuth2AuthorizeParams) {
	if p == nil {
		return nil
	}

	r = &models.OAuth2AuthorizeParams{}
	r.ResponseType = p.ResponseType
	r.ClientID = p.ClientID
	r.Scope = p.Scope
	r.RedirectURI = p.RedirectURI
	r.State = p.State

	return r
}
