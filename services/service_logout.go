package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/restful"
	"time"
)

//可重入
func (s *AccountService) Logout(ctx *restful.Context, accessToken string, refreshToken string) (err error) {
	dbRefreshToken, err := s.accountDB.RefreshToken.GetQuery().RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbRefreshToken == nil {
		return nil
	}

	if dbRefreshToken.IsLogout == 1 {
		return nil
	}

	err = s.accountDB.RefreshToken.GetUpdate().
		IsLogout(1).
		LogoutTime(time.Now().UTC()).
		Update(ctx, nil, dbRefreshToken.Id)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType: models.OperationLogout,
		UserId:        dbRefreshToken.UserId,
	})

	return nil
}
