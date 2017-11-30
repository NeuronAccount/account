package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/storages/oauth"
)

func (s *AccountService) OAuth2AuthorizationCodeToken(code string, redirectUri string, clientId string, oAuth2Client *models.OAuth2Client) (token *models.OAuth2AccessToken, err error) {
	if clientId != oAuth2Client.ClientId {
		return nil, errors.InvalidParam("ClientId", "ClientId不匹配")
	}

	dbAuthorizationCode, err := s.oauthDB.AuthorizationCode.GetQuery().AuthorizationCode_Equal(code).QueryOne(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	if dbAuthorizationCode == nil {
		return nil, errors.InvalidParam("Code", "code不存在")
	}

	if dbAuthorizationCode.ClientId != clientId {
		return nil, errors.InvalidParam("ClientId", "ClientId不匹配")
	}

	if dbAuthorizationCode.RedirectUri != redirectUri {
		return nil, errors.InvalidParam("redirectUri错误", "redirectUri不匹配")
	}

	dbAccessToken := &oauth.AccessToken{}
	dbAccessToken.ClientId = clientId
	dbAccessToken.AccountId = dbAuthorizationCode.AccountId
	dbAccessToken.OauthScope = dbAuthorizationCode.OauthScope
	dbAccessToken.ExpireSeconds = models.OAUTH_ACCESS_TOKEN_EXPIRE_SECONDS
	dbAccessToken.AccessToken = rand.NextBase64(16)
	_, err = s.oauthDB.AccessToken.Insert(context.Background(), nil, dbAccessToken)
	if err != nil {
		return nil, err
	}

	dbRefreshToken := &oauth.RefreshToken{}
	dbRefreshToken.ClientId = clientId
	dbRefreshToken.AccountId = dbAuthorizationCode.AccountId
	dbRefreshToken.OauthScope = dbAuthorizationCode.OauthScope
	dbRefreshToken.RefreshToken = rand.NextBase64(16)
	_, err = s.oauthDB.RefreshToken.Insert(context.Background(), nil, dbRefreshToken)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AccountService) OAuth2RefreshToken(refresh_token string, scope string, oAuth2Client *models.OAuth2Client) (token *models.OAuth2AccessToken, err error) {
	return nil, nil
}

func (s *AccountService) OAuth2Token(request *models.OAuth2TokenRequest, oAuth2Client *models.OAuth2Client) (token *models.OAuth2AccessToken, err error) {
	if request.GrantType == "authorization_code" {
		return s.OAuth2AuthorizationCodeToken(request.AuthorizationCode, request.RedirectURI, request.ClientID, oAuth2Client)
	} else if request.GrantType == "refresh_token" {
		return s.OAuth2RefreshToken(request.RefreshToken, request.Scope, oAuth2Client)
	} else {
		return nil, errors.InvalidParams(&errors.ParamError{Field: "GrantType", Code: "UnknownType", Message: "未知的类型"})
	}

	return nil, nil
}
