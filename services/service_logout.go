package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) Logout(ctx *restful.Context, jwt string) (err error) {
	s.addOperation(ctx, &models.Operation{
		OperationType: models.OperationLogout,
		AccountID:     "", //todo jwt->accountID
	})

	return nil
}
