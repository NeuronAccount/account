package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) Logout(ctx *rest.Context, userId string) (err error) {
	_, err = s.accountDB.RefreshToken.Query().UserIdEqual(userId).Delete(ctx, nil)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType: models.OperationLogout,
		UserId:        userId,
	})

	return nil
}
