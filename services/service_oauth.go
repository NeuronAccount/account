package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) OauthState(ctx *restful.Context) (state string, err error) {
	return "", nil
}

func (s *AccountService) OauthJump(ctx *restful.Context, params *models.OauthJumpParams) (
	userToken *models.UserToken, err error) {
	return nil, nil
}
