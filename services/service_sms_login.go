package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) SmsLogin(ctx *restful.Context, phone string, smsCode string) (jwt string, err error) {
	//check account exists
	dbAccount, err := s.accountDB.Account.GetQuery().
		PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}
	if dbAccount == nil {
		return "", errors.NotFound("帐号不存在")
	}

	err = s.validateSmsCode(ctx, models.SmsSceneLogin, phone, smsCode)
	if err != nil {
		return "", err
	}

	//gen gwt
	jwt, err = s.generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	s.addOperation(ctx, &models.Operation{
		OperationType: models.OperationSmsLogin,
		Phone:         phone,
		AccountID:     dbAccount.AccountId,
	})

	return jwt, nil
}
