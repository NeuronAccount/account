package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) ResetPassword(ctx *restful.Context, phone string, smsCode string, newPassword string) (err error) {
	err = s.validateSmsCode(ctx, models.SmsSceneResetPassword, phone, smsCode)
	if err != nil {
		return err
	}

	dbAccount, err := s.accountDB.Account.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return nil
	}

	if dbAccount == nil {
		return errors.NotFound("帐号不存在")
	}

	dbAccount.PasswordHash = s.calcPasswordHash(newPassword)
	err = s.accountDB.Account.Update(ctx, nil, dbAccount)
	if err != nil {
		return err
	}

	s.addOperation(nil, &models.Operation{
		OperationType: models.OperationResetPassword,
		Phone:         phone,
		AccountID:     dbAccount.AccountId,
	})

	return nil
}
