package oauth

import "github.com/NeuronGroup/Account/models"

func FromOauthClient(p *OauthClient) (r *models.OAuth2Client) {
	if p == nil {
		return nil
	}

	r = &models.OAuth2Client{}
	r.ClientId = p.ClientId
	r.PasswordHash = p.PasswordHash
	r.AccountId = p.AccountId
	r.RedirectUri = p.RedirectUri

	return r
}
