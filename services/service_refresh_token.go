package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) RefreshToken(ctx *rest.Context, refreshToken string) (userToken *models.UserToken, err error) {
	dbRefreshToken, err := s.accountDB.RefreshToken.Query().RefreshTokenEqual(refreshToken).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshToken == nil {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}

	return s.createUserToken(ctx, dbRefreshToken.UserId)
}
