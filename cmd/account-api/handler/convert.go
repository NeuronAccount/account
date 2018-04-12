package handler

import "github.com/NeuronAccount/account/models"
import api "github.com/NeuronAccount/account/api/gen/models"

func fromToken(p *models.UserToken) (r *api.UserToken) {
	if p == nil {
		return nil
	}

	r = &api.UserToken{}
	r.AccessToken = &p.AccessToken
	r.RefreshToken = &p.RefreshToken

	return r
}
