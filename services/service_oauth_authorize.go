package services

import (
	"context"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/storages/oauth"
)

func (s *AccountService) OauthAuthorize(jwt string, params *models.OAuth2AuthorizeParams) (code *models.OauthAuthorizationCode, err error) {
	tx, err := s.oauthDB.BeginReadCommittedTx(context.Background(), false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	accountId := jwt

	dbAuthorizeCode, err := s.oauthDB.AuthorizationCode.GetQuery().ForUpdate().
		ClientId_Equal(params.ClientID).And().AccountId_Equal(accountId).
		QueryOne(context.Background(), nil)

	if dbAuthorizeCode == nil {
		dbAuthorizeCode = &oauth.AuthorizationCode{}
		dbAuthorizeCode.ClientId = params.ClientID
		dbAuthorizeCode.AccountId = accountId
		dbAuthorizeCode.RedirectUri = params.RedirectURI
		dbAuthorizeCode.OauthScope = params.Scope
		dbAuthorizeCode.ExpireSeconds = 300
		_, err = s.oauthDB.AuthorizationCode.Insert(context.Background(), tx, dbAuthorizeCode)
		if err != nil {
			return nil, err
		}
	} else {
		dbAuthorizeCode.ClientId = params.ClientID
		dbAuthorizeCode.AccountId = accountId
		dbAuthorizeCode.RedirectUri = params.RedirectURI
		dbAuthorizeCode.OauthScope = params.Scope
		dbAuthorizeCode.ExpireSeconds = 300
		err = s.oauthDB.AuthorizationCode.Update(context.Background(), tx, dbAuthorizeCode)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
