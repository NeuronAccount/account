package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) OauthState(ctx *rest.Context) (state string, err error) {
	return "", nil
}

func (s *AccountService) OauthJump(ctx *rest.Context, params *models.OauthJumpParams) (
	userToken *models.UserToken, err error) {
	return nil, nil
}

func (s *AccountService) BindOauthAccount(ctx *rest.Context, userId string) (err error) {
	return nil
}

func (s *AccountService) UnbindOauthAccount(ctx *rest.Context, userId string) (err error) {
	return nil
}
