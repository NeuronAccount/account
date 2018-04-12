package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/restful"
	"time"
)

func (s *AccountService) Logout(ctx *restful.Context, accessToken string, refreshToken string) (err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbRefreshToken == nil {
		return nil
	}

	if dbRefreshToken.IsLogout == 1 {
		return nil
	}

	dbRefreshToken.IsLogout = 1
	dbRefreshToken.LogoutTime = time.Now().UTC()
	err = s.userDB.RefreshToken.Update(ctx, nil, dbRefreshToken)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType: models.OperationLogout,
		UserId:        dbRefreshToken.UserId,
	})

	return nil
}
