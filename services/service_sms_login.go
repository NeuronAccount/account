package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
)

func (s *AccountService) SmsLogin(ctx context.Context, phone string, smsCode string) (jwt string, err error) {
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

	return jwt, nil
}
