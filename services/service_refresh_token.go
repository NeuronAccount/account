package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) RefreshToken(ctx *rest.Context, refreshToken string) (userToken *models.UserToken, err error) {
	//检查RefreshToken是否最新且有效
	dbRefreshToken, err := s.accountDB.RefreshToken.GetQuery().
		OrderBy(neuron_account_db.REFRESH_TOKEN_FIELD_ID, false).
		Limit(0, 1).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshToken == nil {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}
	if dbRefreshToken.IsLogout == 1 {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}
	if dbRefreshToken.RefreshToken != refreshToken {
		dbRefreshTokenOld, err := s.accountDB.RefreshToken.GetQuery().
			RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
		if err != nil {
			return nil, err
		}

		if dbRefreshTokenOld != nil {
			return nil, errors.NotFound("您已在其它地方登录，请重新登录")
		}

		return nil, errors.NotFound("Token已失效，请重新登录")
	}

	return s.createUserToken(ctx, dbRefreshToken.UserId, dbRefreshToken)
}
